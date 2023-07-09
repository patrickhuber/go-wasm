package io

import (
	"bytes"
	"encoding/binary"
	"math"

	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"

	"github.com/patrickhuber/go-wasm/abi/trap"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
)

func Load(cx *types.CallContext, t types.ValType, ptr uint32) (res Result[any]) {
	defer handle.Error(&res)

	t = Despecialize(t)

	switch vt := t.(type) {
	case types.Bool:
		b := LoadBool(cx, ptr).Unwrap()
		return Cast[bool, any](b)
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
		i := LoadUInt32(cx, ptr).Unwrap()
		return LiftOwn(cx, i, vt)
	case types.Borrow:
		i := LoadUInt32(cx, ptr).Unwrap()
		return LiftBorrow(cx, i, vt)
	}
	return result.Errorf[any]("unrecognized type %T", t)
}

func LoadChar(cx *types.CallContext, ptr uint32, t types.ValType) (res Result[any]) {
	defer handle.Error(&res)
	t = Castf[types.ValType, types.Char](t, "load char unable to match type").Unwrap()
	i := LoadIntWithSize(cx, ptr, SizeOfChar, false).Unwrap()
	u32 := Cast[any, uint32](i).Unwrap()
	r := ConvertU32ToRune(u32).Unwrap()
	return result.Ok[any](r)
}

func ConvertU32ToRune(u32 uint32) (res Result[rune]) {
	defer handle.Error(&res)
	trap.Iff(u32 >= 0x110000, "u32 %d >= 0x110000", u32)
	trap.Iff(0xd800 <= u32 && u32 <= 0xdfff, " 0xd800 <= %d <= 0xdfff", u32)
	return result.Ok(rune(u32))
}

func LoadBool(cx *types.CallContext, ptr uint32) (res Result[bool]) {
	defer handle.Error(&res)
	i := LoadInt(cx, ptr, types.NewU8()).Unwrap()
	u8 := Cast[any, uint8](i).Unwrap()
	return result.Ok(u8 != 0)
}

func LoadUInt32(cx *types.CallContext, ptr uint32) (res Result[uint32]) {
	defer handle.Error(&res)
	val := LoadInt(cx, ptr, types.NewU32()).Unwrap()
	return Cast[any, uint32](val)
}

func LoadUInt64(cx *types.CallContext, ptr uint32) (res Result[uint64]) {
	defer handle.Error(&res)
	val := LoadInt(cx, ptr, types.NewU64()).Unwrap()
	return Cast[any, uint64](val)
}

func LoadInt(c *types.CallContext, ptr uint32, t types.ValType) (res Result[any]) {
	defer handle.Error(&res)

	size := Size(t).Unwrap()
	buf := c.Options.Memory.Bytes()[ptr : ptr+size]

	var value any
	switch t.(type) {
	case types.U8:
		value = buf[0]
	case types.U16:
		value = binary.LittleEndian.Uint16(buf)
	case types.U32:
		value = binary.LittleEndian.Uint32(buf)
	case types.U64:
		value = binary.LittleEndian.Uint64(buf)
	case types.S8:
		value = int8(buf[0])
	case types.S16:
		value = int16(binary.LittleEndian.Uint16(buf))
	case types.S32:
		value = int32(binary.LittleEndian.Uint32(buf))
	case types.S64:
		value = int64(binary.LittleEndian.Uint64(buf))
	default:
		value = uint32(0)
	}
	return result.Ok(value)
}

func LoadIntWithSize(c *types.CallContext, ptr uint32, nbytes uint32, sign bool) Result[any] {
	var t types.ValType
	switch {
	case nbytes == 0:
		return result.Ok[any](uint32(0))
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

func LoadFloat(cx *types.CallContext, ptr uint32, t types.ValType) (res Result[any]) {
	defer handle.Error(&res)
	switch t.(type) {
	case types.Float32:
		i := LoadUInt32(cx, ptr).Unwrap()
		f := math.Float32frombits(i)
		return result.Ok[any](f)
	case types.Float64:
		i := LoadUInt64(cx, ptr).Unwrap()
		f := math.Float64frombits(i)
		return result.Ok[any](f)
	}
	return result.Errorf[any]("LoadFloat: invalid float type %T", t)
}

func LoadString(cx *types.CallContext, ptr uint32) (res Result[string]) {
	defer handle.Error(&res)
	begin := LoadUInt32(cx, ptr).Unwrap()

	// is this byte order mark?
	taggedCodeUnits := LoadUInt32(cx, ptr+4).Unwrap()
	return LoadStringFromRange(cx, begin, taggedCodeUnits)
}

func LoadStringFromRange(cx *types.CallContext, ptr, taggedCodeUnits uint32) (res Result[string]) {
	defer handle.Error(&res)
	
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
	align := AlignTo(ptr, uint32(codec.Alignment()))

	trap.Iff(ptr != align, "error aligning ptr %d to %d", ptr, uint32(codec.Alignment()))
	trap.Iff(ptr+byteLength > uint32(cx.Options.Memory.Len(),"destination %d > memory size %d", ptr+byteLength, cx.Options.Memory.Len())
	
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
