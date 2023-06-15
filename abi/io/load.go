package io

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
)

func Load(cx *types.Context, t types.ValType, ptr uint32) (any, error) {
	k := t.Kind()
	switch k {
	case kind.Bool:
		return LoadBool(cx, ptr)
	case kind.U8:
		fallthrough
	case kind.U16:
		fallthrough
	case kind.U32:
		fallthrough
	case kind.U64:
		size, err := t.Size()
		if err != nil {
			return nil, err
		}
		return LoadInt(cx, ptr, size, false)
	case kind.S8:
		fallthrough
	case kind.S16:
		fallthrough
	case kind.S32:
		fallthrough
	case kind.S64:
		size, err := t.Size()
		if err != nil {
			return nil, err
		}
		return LoadInt(cx, ptr, size, true)
	case kind.Float32:
		fallthrough
	case kind.Float64:
		size, err := t.Size()
		if err != nil {
			return nil, err
		}
		return LoadFloat(cx, ptr, size)
	case kind.Char:
		size, err := t.Size()
		if err != nil {
			return nil, err
		}
		return LoadChar(cx, ptr, size)
	case kind.String:
		return LoadString(cx, ptr)
	case kind.List:
		l := t.(*types.List)
		return LoadList(cx, ptr, l.Type)
	case kind.Record:
		r := t.(*types.Record)
		return LoadRecord(cx, ptr, r.Fields)
	case kind.Variant:
		v := t.(*types.Variant)
		return LoadVariant(cx, ptr, v)
	case kind.Flags:
		f := t.(*types.Flags)
		return LoadFlags(cx, ptr, f)
	case kind.Own:
		i, err := LoadUInt32(cx, ptr)
		if err != nil {
			return nil, err
		}
		o := t.(*types.Own)
		return LiftOwn(cx, i, o)
	case kind.Borrow:
		i, err := LoadUInt32(cx, ptr)
		if err != nil {
			return nil, err
		}
		b := t.(*types.Borrow)
		return LiftBorrow(cx, i, b)
	}
	return nil, fmt.Errorf("unrecognized type %s", k.String())
}

func LoadChar(cx *types.Context, ptr uint32, nbytes uint32) (any, error) {
	i, err := LoadInt(cx, ptr, nbytes, false)
	if err != nil {
		return nil, err
	}
	return rune(i.(uint32)), nil
}

func LoadBool(cx *types.Context, ptr uint32) (bool, error) {
	i, err := LoadInt(cx, ptr, 1, false)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

func LoadUInt32(cx *types.Context, ptr uint32) (uint32, error) {
	val, err := LoadInt(cx, ptr, 4, false)
	if err != nil {
		return 0, err
	}
	return val.(uint32), nil
}

func LoadInt(c *types.Context, ptr uint32, nbytes uint32, signed bool) (any, error) {
	buf := c.Options.Memory.Bytes()[ptr : ptr+nbytes]
	switch nbytes {
	case 1:
		if signed {
			return int8(buf[0]), nil
		}
		return buf[0], nil
	case 2:
		v := binary.LittleEndian.Uint16(buf)
		if signed {
			return int16(v), nil
		}
		return v, nil
	case 4:
		v := binary.LittleEndian.Uint32(buf)
		if signed {
			return int32(v), nil
		}
		return v, nil
	case 8:
		v := binary.LittleEndian.Uint64(buf)
		if signed {
			return int64(v), nil
		}
		return v, nil
	}
	return nil, fmt.Errorf("invalid type")
}

func LoadFloat(cx *types.Context, ptr uint32, nbytes uint32) (any, error) {
	i, err := LoadInt(cx, ptr, nbytes, false)
	if err != nil {
		return nil, err
	}
	if nbytes == 4 {
		ui := i.(uint32)
		return math.Float32frombits(ui), nil
	}
	ui := i.(uint64)
	return math.Float64frombits(ui), nil
}

func LoadString(cx *types.Context, ptr uint32) (string, error) {
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

func LoadStringFromRange(cx *types.Context, ptr, taggedCodeUnits uint32) (string, error) {

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

	err = types.TrapIf(ptr != types.AlignTo(ptr, uint32(codec.Alignment())))
	if err != nil {
		return "", err
	}

	err = types.TrapIf(ptr+byteLength > uint32(cx.Options.Memory.Len()))
	if err != nil {
		return "", err
	}

	buf := cx.Options.Memory.Bytes()[ptr : ptr+byteLength]
	return encoding.DecodeString(codec, bytes.NewReader(buf))
}

func LoadList(cx *types.Context, ptr uint32, elementType types.ValType) ([]any, error) {

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

func LoadListFromRange(cx *types.Context, ptr uint32, length uint32, elementType types.ValType) ([]any, error) {
	alignment, err := elementType.Alignment()
	if err != nil {
		return nil, err
	}
	err = types.TrapIf(ptr != types.AlignTo(ptr, alignment))
	if err != nil {
		return nil, err
	}
	size, err := elementType.Size()
	if err != nil {
		return nil, err
	}
	err = types.TrapIf(ptr+length*size > uint32(cx.Options.Memory.Len()))
	if err != nil {
		return nil, err
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

func LoadRecord(cx *types.Context, ptr uint32, fields []types.Field) (map[string]any, error) {
	record := map[string]any{}
	for _, field := range fields {
		alignment, err := field.Type.Alignment()
		if err != nil {
			return nil, err
		}
		ptr = types.AlignTo(ptr, alignment)
		val, err := Load(cx, field.Type, ptr)
		if err != nil {
			return nil, err
		}
		record[field.Label] = val
		size, err := field.Type.Size()
		if err != nil {
			return nil, err
		}
		ptr += size
	}
	return record, nil
}

// LoadVariant loads the variant from the context at the ptr
func LoadVariant(cx *types.Context, ptr uint32, v *types.Variant) (map[string]any, error) {
	dt, err := v.DiscriminantType()
	if err != nil {
		return nil, err
	}
	discSize, err := dt.Size()
	if err != nil {
		return nil, err
	}
	caseIndex, err := LoadInt(cx, ptr, discSize, false)
	var u32CaseIndex uint32 = 0
	switch dt.Kind() {
	case kind.U8:
		u32CaseIndex = uint32(caseIndex.(uint8))
	case kind.U16:
		u32CaseIndex = uint32(caseIndex.(uint16))
	case kind.U32:
		u32CaseIndex = caseIndex.(uint32)
	case kind.U64:
		// could cause problems if caseIndex is actually a u64
		u32CaseIndex = uint32(caseIndex.(uint64))
	}
	if err != nil {
		return nil, err
	}
	ptr += discSize
	err = types.TrapIf(u32CaseIndex >= uint32(len(v.Cases)))
	if err != nil {
		return nil, err
	}
	c := v.Cases[u32CaseIndex]
	maxCaseAlignment, err := v.MaxCaseAlignment()
	if err != nil {
		return nil, err
	}
	ptr = types.AlignTo(ptr, maxCaseAlignment)

	caseLabel := v.CaseLabelWithRefinements(c)
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

func LoadFlags(cx *types.Context, ptr uint32, flags *types.Flags) (map[string]bool, error) {
	size, err := flags.Size()
	if err != nil {
		return nil, err
	}

	i, err := LoadInt(cx, ptr, size, false)
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
	flagMap := map[string]bool{}
	for _, label := range flags.Labels {
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

func UnpackFlagsFromInt(i int, labels []string) map[string]any {
	unpacked := map[string]any{}
	for _, label := range labels {
		unpacked[label] = (i&1 == 1)
		i >>= 1
	}
	return unpacked
}
