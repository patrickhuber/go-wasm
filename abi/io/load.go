package io

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
)

func Load(cx *types.CallContext, t types.ValType, ptr uint32) (any, error) {
	t = Despecialize(t)

	switch vt := t.(type) {
	case types.Bool:
		return LoadBool(cx, ptr)
	case types.U8:
		return LoadInt(cx, ptr, t)
	case types.U16:
		return LoadInt(cx, ptr, t)
	case types.U32:
		return LoadInt(cx, ptr, t)
	case types.U64:
		return LoadInt(cx, ptr, t)
	case types.S8:
		return LoadInt(cx, ptr, t)
	case types.S16:
		return LoadInt(cx, ptr, t)
	case types.S32:
		return LoadInt(cx, ptr, t)
	case types.S64:
		return LoadInt(cx, ptr, t)
	case types.Float32:
		return LoadFloat(cx, ptr, t)
	case types.Float64:
		return LoadFloat(cx, ptr, t)
	case types.Char:
		return LoadChar(cx, ptr, t)
	case types.String:
		return LoadString(cx, ptr)
	case types.List:
		return LoadList(cx, ptr, vt.Type())
	case types.Record:
		return LoadRecord(cx, ptr, vt.Fields())
	case types.Variant:
		return LoadVariant(cx, ptr, vt)
	case types.Flags:
		return LoadFlags(cx, ptr, vt)
	case types.Own:
		i, err := LoadUInt32(cx, ptr)
		if err != nil {
			return nil, err
		}
		return LiftOwn(cx, i, vt)
	case types.Borrow:
		i, err := LoadUInt32(cx, ptr)
		if err != nil {
			return nil, err
		}
		return LiftBorrow(cx, i, vt)
	}
	return nil, fmt.Errorf("unrecognized type %T", t)
}

func LoadChar(cx *types.CallContext, ptr uint32, t types.ValType) (any, error) {
	switch t.(type) {
	case types.Char:
	default:
		return nil, fmt.Errorf("LoadChar unable to match type %T", t)
	}

	i, err := LoadIntWithSize(cx, ptr, SizeOfChar, false)
	if err != nil {
		return nil, err
	}

	u32, ok := i.(uint32)
	if !ok {
		return nil, types.NewCastError(i, "uint32")
	}

	return ConvertU32ToRune(u32)
}

func ConvertU32ToRune(u32 uint32) (rune, error) {
	if u32 >= 0x110000 {
		return 0, types.TrapWith("u32 %d >= 0x110000", u32)
	}
	if 0xd800 <= u32 && u32 <= 0xdfff {
		return 0, types.TrapWith(" 0xd800 <= %d <= 0xdfff", u32)
	}
	return rune(u32), nil
}

func LoadBool(cx *types.CallContext, ptr uint32) (bool, error) {
	i, err := LoadInt(cx, ptr, types.NewU8())
	if err != nil {
		return false, err
	}
	u8, ok := i.(uint8)
	if !ok {
		return false, types.NewCastError(i, "uint8")
	}
	return u8 != 0, nil
}

func LoadUInt32(cx *types.CallContext, ptr uint32) (uint32, error) {
	val, err := LoadInt(cx, ptr, types.NewU32())
	if err != nil {
		return 0, err
	}
	u32, ok := val.(uint32)
	if !ok {
		return 0, types.NewCastError(val, "uint32")
	}
	return u32, nil
}

func LoadUInt64(cx *types.CallContext, ptr uint32) (uint64, error) {
	val, err := LoadInt(cx, ptr, types.NewU64())
	if err != nil {
		return 0, err
	}
	u64, ok := val.(uint64)
	if !ok {
		return 0, types.NewCastError(val, "uint64")
	}
	return u64, nil
}

func LoadInt(c *types.CallContext, ptr uint32, t types.ValType) (any, error) {
	size, err := Size(t)
	if err != nil {
		return nil, err
	}
	buf := c.Options.Memory.Bytes()[ptr : ptr+size]
	switch t.(type) {
	case types.U8:
		return buf[0], nil
	case types.U16:
		return binary.LittleEndian.Uint16(buf), nil
	case types.U32:
		return binary.LittleEndian.Uint32(buf), nil
	case types.U64:
		return binary.LittleEndian.Uint64(buf), nil
	case types.S8:
		return int8(buf[0]), nil
	case types.S16:
		return int16(binary.LittleEndian.Uint16(buf)), nil
	case types.S32:
		return int32(binary.LittleEndian.Uint32(buf)), nil
	case types.S64:
		return int64(binary.LittleEndian.Uint64(buf)), nil
	default:
		return uint32(0), nil
	}
}

func LoadIntWithSize(c *types.CallContext, ptr uint32, nbytes uint32, sign bool) (any, error) {
	var t types.ValType
	switch {
	case nbytes == 0:
		return uint32(0), nil
	case nbytes == 1 && !sign:
		t = types.NewU8()
	case nbytes == 2 && !sign:
		t = types.NewU16()
	case nbytes == 4 && !sign:
		t = types.NewU32()
	case nbytes == 8 && !sign:
		t = types.NewU64()
	case nbytes == 1 && sign:
		t = types.NewS8()
	case nbytes == 2 && sign:
		t = types.NewS16()
	case nbytes == 4 && sign:
		t = types.NewS32()
	case nbytes == 8 && sign:
		t = types.NewS64()
	}
	return LoadInt(c, ptr, t)
}

func LoadFloat(cx *types.CallContext, ptr uint32, t types.ValType) (any, error) {
	switch t.(type) {
	case types.Float32:
		i, err := LoadUInt32(cx, ptr)
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(i), nil
	case types.Float64:
		i, err := LoadUInt64(cx, ptr)
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(i), nil
	}
	return nil, fmt.Errorf("LoadFloat: invalid float type %T", t)
}

func LoadString(cx *types.CallContext, ptr uint32) (string, error) {
	begin, err := LoadUInt32(cx, ptr)
	if err != nil {
		return "", err
	}
	// is this byte order mark?
	taggedCodeUnits, err := LoadUInt32(cx, ptr+4)
	if err != nil {
		return "", err
	}
	return LoadStringFromRange(cx, begin, taggedCodeUnits)
}

func LoadStringFromRange(cx *types.CallContext, ptr, taggedCodeUnits uint32) (string, error) {

	srcEncoding := cx.Options.StringEncoding
	tcu := UInt32ToTaggedCodeUnits(taggedCodeUnits)
	if srcEncoding == encoding.Latin1Utf16 {
		if tcu.UTF16 {
			srcEncoding = encoding.UTF16LE
		} else {
			srcEncoding = encoding.Latin1
		}
	}

	codec, err := encoding.DefaultFactory().Get(srcEncoding)
	if err != nil {
		return "", err
	}

	byteLength := tcu.CodeUnits * uint32(codec.RuneSize())
	align, err := AlignTo(ptr, uint32(codec.Alignment()))

	if err != nil {
		return "", err
	}

	if ptr != align {
		return "", types.TrapWith("error aligning ptr %d to %d", ptr, uint32(codec.Alignment()))
	}

	if ptr+byteLength > uint32(cx.Options.Memory.Len()) {
		return "", types.TrapWith("destination %d > memory size %d", ptr+byteLength, cx.Options.Memory.Len())
	}

	buf := cx.Options.Memory.Bytes()[ptr : ptr+byteLength]
	return encoding.DecodeString(codec, bytes.NewReader(buf))
}

func LoadList(cx *types.CallContext, ptr uint32, elementType types.ValType) ([]any, error) {

	begin, err := LoadUInt32(cx, ptr)
	if err != nil {
		return nil, err
	}
	length, err := LoadUInt32(cx, ptr+4)
	if err != nil {
		return nil, err
	}

	return LoadListFromRange(cx, begin, length, elementType)
}

func LoadListFromRange(cx *types.CallContext, ptr uint32, length uint32, elementType types.ValType) ([]any, error) {
	alignment, err := Alignment(elementType)
	if err != nil {
		return nil, err
	}
	align, err := AlignTo(ptr, alignment)
	if err != nil {
		return nil, err
	}
	if ptr != align {
		return nil, types.TrapWith("unable to align ptr %d with %d", ptr, alignment)
	}

	size, err := Size(elementType)
	if err != nil {
		return nil, err
	}
	if ptr+length*size > uint32(cx.Options.Memory.Len()) {
		return nil, types.TrapWith("destination size %d is greater than memory size %d", ptr+length*size, cx.Options.Memory.Len())
	}
	var list []any
	var i uint32 = 0
	for ; i < length; i++ {
		element, err := Load(cx, elementType, ptr+i*size)
		if err != nil {
			return nil, err
		}
		list = append(list, element)
	}
	return list, nil
}

func LoadRecord(cx *types.CallContext, ptr uint32, fields []types.Field) (map[string]any, error) {
	record := map[string]any{}
	for _, field := range fields {
		alignment, err := Alignment(field.Type)
		if err != nil {
			return nil, err
		}
		ptr, err = AlignTo(ptr, alignment)
		if err != nil {
			return nil, err
		}
		val, err := Load(cx, field.Type, ptr)
		if err != nil {
			return nil, err
		}
		record[field.Label] = val
		size, err := Size(field.Type)
		if err != nil {
			return nil, err
		}
		ptr += size
	}
	return record, nil
}

// LoadVariant loads the variant from the context at the ptr
func LoadVariant(cx *types.CallContext, ptr uint32, v types.Variant) (map[string]any, error) {
	dt, err := DiscriminantType(v.Cases())
	if err != nil {
		return nil, err
	}
	caseIndex, err := LoadInt(cx, ptr, dt)
	var u32CaseIndex uint32 = 0
	switch dt.(type) {
	case types.U8:
		u32CaseIndex = uint32(caseIndex.(uint8))
	case types.U16:
		u32CaseIndex = uint32(caseIndex.(uint16))
	case types.U32:
		u32CaseIndex = caseIndex.(uint32)
	case types.U64:
		// could cause problems if caseIndex is actually a u64
		u32CaseIndex = uint32(caseIndex.(uint64))
	}
	if err != nil {
		return nil, err
	}

	discSize, err := Size(dt)
	if err != nil {
		return nil, err
	}

	ptr += discSize
	if u32CaseIndex >= uint32(len(v.Cases())) {
		return nil, types.TrapWith("case index %d is outside the bounds of the case index length %d", u32CaseIndex, len(v.Cases()))
	}

	c := v.Cases()[u32CaseIndex]
	maxCaseAlignment, err := MaxCaseAlignment(v.Cases())
	if err != nil {
		return nil, err
	}

	ptr, err = AlignTo(ptr, maxCaseAlignment)
	if err != nil {
		return nil, err
	}

	caseLabel := CaseLabelWithRefinements(c, v.Cases())
	var value any
	if c.Type == nil {
		value = nil
	} else {
		value, err = Load(cx, c.Type, ptr)
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{
		caseLabel: value,
	}, nil
}

func CaseLabelWithRefinements(c types.Case, cases []types.Case) string {
	label := c.Label
	for {
		if c.Refines == nil {
			break
		}
		index := FindCaseIndex(cases, *c.Refines)
		if index < 0 {
			break
		}
		c = cases[index]
		label += "|" + c.Label
	}
	return label
}

func FindCaseIndex(cases []types.Case, label string) int {
	for i, c := range cases {
		if c.Label == label {
			return i
		}
	}
	return -1
}

func LoadFlags(cx *types.CallContext, ptr uint32, flags types.Flags) (map[string]any, error) {
	size, err := Size(flags)
	if err != nil {
		return nil, err
	}

	i, err := LoadIntWithSize(cx, ptr, size, false)
	if err != nil {
		return nil, err
	}

	var ui uint64 = 0
	switch size {
	case 1:
		ui = uint64(i.(uint8))
	case 2:
		ui = uint64(i.(uint16))
	case 4:
		ui = uint64(i.(uint32))
	case 8:
		ui = i.(uint64)
	}
	flagMap := map[string]any{}
	for _, label := range flags.Labels() {
		v := ui & 1
		b := false
		if v > 0 {
			b = true
		}
		flagMap[label] = b
		ui >>= 1
	}
	return flagMap, nil
}

func UnpackFlagsFromInt(i uint64, labels []string) map[string]any {
	unpacked := map[string]any{}
	for _, label := range labels {
		unpacked[label] = (i&1 == 1)
		i >>= 1
	}
	return unpacked
}
