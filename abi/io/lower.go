package io

import (
	"fmt"
	"math"
	"strings"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func LowerFlat(cx *types.CallContext, v any, t types.ValType) ([]values.Value, error) {
	t = t.Despecialize()
	k := t.Kind()
	switch k {
	case kind.Bool:
		return LowerBool(v)
	case kind.S8:
		return LowerS8(v)
	case kind.S16:
		return LowerS16(v)
	case kind.S32:
		return LowerS32(v)
	case kind.S64:
		return LowerS64(v)
	case kind.U8:
		return LowerU8(v)
	case kind.U16:
		return LowerU16(v)
	case kind.U32:
		return LowerU32(v)
	case kind.U64:
		return LowerU64(v)
	case kind.Float32:
		return LowerFloat32(v)
	case kind.Float64:
		return LowerFloat64(v)
	case kind.Char:
		return LowerChar(v)
	case kind.String:
		return LowerString(cx, v)
	case kind.List:
		l, ok := t.(*types.List)
		if !ok {
			return nil, types.NewCastError(t, "*types.Record")
		}
		return LowerFlatList(cx, v, l.Type)
	case kind.Record:
		r, ok := t.(*types.Record)
		if !ok {
			return nil, types.NewCastError(t, "*types.Record")
		}
		return LowerFlatRecord(cx, v, r)
	case kind.Flags:
		f, ok := t.(*types.Flags)
		if !ok {
			return nil, types.NewCastError(t, "*types.Flags")
		}
		return LowerFlatFlags(cx, v, f)
	case kind.Variant:
		variant, ok := t.(*types.Variant)
		if !ok {
			return nil, types.NewCastError(t, "*types.Variant")
		}
		return LowerFlatVariant(cx, v, variant)
	}
	return nil, fmt.Errorf("unable to lower type %s", k.String())
}

func LowerBool(v any) ([]values.Value, error) {
	var i values.U32
	b, ok := v.(bool)
	if !ok {
		return nil, types.NewCastError(v, "bool")
	}
	if b {
		i = 1
	} else {
		i = 0
	}
	return slice(i), nil
}

func LowerU8(v any) ([]values.Value, error) {
	u8, ok := v.(uint8)
	if !ok {
		return nil, types.NewCastError(v, "uint8")
	}
	i := values.U32(u8)
	return slice(i), nil
}

func LowerU16(v any) ([]values.Value, error) {
	u16, ok := v.(uint16)
	if !ok {
		return nil, types.NewCastError(v, "uint16")
	}
	i := values.U32(u16)
	return slice(i), nil
}

func LowerU32(v any) ([]values.Value, error) {
	u32, ok := v.(uint32)
	if !ok {
		return nil, types.NewCastError(v, "uint32")
	}
	i := values.U32(u32)
	return slice(i), nil
}

func LowerU64(v any) ([]values.Value, error) {
	u64, ok := v.(uint64)
	if !ok {
		return nil, types.NewCastError(v, "uint64")
	}
	i := values.U64(u64)
	return slice(i), nil
}

func LowerS8(v any) ([]values.Value, error) {
	s8, ok := v.(int8)
	if !ok {
		return nil, types.NewCastError(v, "int8")
	}
	i := values.U32(s8)
	return slice(i), nil
}

func LowerS16(v any) ([]values.Value, error) {
	s16, ok := v.(int16)
	if !ok {
		return nil, types.NewCastError(v, "int16")
	}
	i := values.U32(s16)
	return slice(i), nil
}

func LowerS32(v any) ([]values.Value, error) {
	s32, ok := v.(int32)
	if !ok {
		return nil, types.NewCastError(v, "int32")
	}
	i := values.U32(s32)
	return slice(i), nil
}

func LowerS64(v any) ([]values.Value, error) {
	s64, ok := v.(int64)
	if !ok {
		return nil, types.NewCastError(v, "int64")
	}
	i := values.U64(s64)
	return slice(i), nil
}

func LowerFloat32(v any) ([]values.Value, error) {
	f32, ok := v.(float32)
	if !ok {
		return nil, types.NewCastError(v, "float32")
	}
	f := values.Float32(f32)
	return slice(f), nil
}

func LowerFloat64(v any) ([]values.Value, error) {
	f64, ok := v.(float64)
	if !ok {
		return nil, types.NewCastError(v, "float64")
	}
	f := values.Float64(f64)
	return slice(f), nil
}

func LowerChar(v any) ([]values.Value, error) {
	r := v.(rune)
	i := values.U32(r)
	return slice(i), nil
}

func LowerString(cx *types.CallContext, v any) ([]values.Value, error) {
	str := v.(string)
	ptr, packedLength, err := StoreStringIntoRange(cx, str)
	if err != nil {
		return nil, err
	}
	iptr, err := LowerU32(ptr)
	if err != nil {
		return nil, err
	}
	ilen, err := LowerU32(packedLength)
	if err != nil {
		return nil, err
	}
	return append(iptr, ilen...), nil
}

func LowerFlatList(cx *types.CallContext, v any, t types.ValType) ([]values.Value, error) {
	ptr, length, err := StoreListIntoRange(cx, v, t)
	if err != nil {
		return nil, err
	}
	return []values.Value{
		values.U32(ptr),
		values.U32(length),
	}, nil
}

func LowerFlatRecord(cx *types.CallContext, v any, r *types.Record) ([]values.Value, error) {
	var flat []values.Value
	vMap, ok := v.(map[string]any)
	if !ok {
		return nil, types.NewCastError(v, "map[string]any")
	}
	for _, field := range r.Fields {
		lowerFields, err := LowerFlat(cx, vMap[field.Label], field.Type)
		if err != nil {
			return nil, err
		}
		flat = append(flat, lowerFields...)
	}
	return flat, nil
}

func LowerFlatFlags(cx *types.CallContext, v any, f *types.Flags) ([]values.Value, error) {
	vMap, ok := v.(map[string]any)
	if !ok {
		return nil, types.NewCastError(v, "map[string]any")
	}
	packed, err := PackFlagsIntoInt(vMap, f)
	if err != nil {
		return nil, err
	}
	var flat []values.Value
	numFlags := f.NumI32Flags()
	for i := 0; i < int(numFlags); i++ {
		u32 := values.U32(packed & 0xffffffff)
		flat = append(flat, u32)
		packed >>= 32
	}
	if packed != 0 {
		return nil, fmt.Errorf("invalid flag value")
	}
	return flat, nil
}

func LowerFlatVariant(cx *types.CallContext, v any, variant *types.Variant) ([]values.Value, error) {
	caseIndex, caseValue, err := MatchCase(v, variant.Cases)
	if err != nil {
		return nil, err
	}
	flatTypes, err := variant.Flatten()
	if err != nil {
		return nil, err
	}
	if len(flatTypes) == 0 {
		return nil, fmt.Errorf("expected at least one flattend type")
	}
	first := flatTypes[0]
	flatTypes = flatTypes[1:]
	if first != kind.U32 {
		return nil, fmt.Errorf("expected kind.U32")
	}
	c := variant.Cases[caseIndex]
	var payload []values.Value
	if c.Type == nil {
		payload = nil
	} else {
		payload, err = LowerFlat(cx, caseValue, c.Type)
		if err != nil {
			return nil, err
		}
	}
	for i, have := range payload {
		if len(flatTypes) == 0 {
			return nil, fmt.Errorf("expected len flatTypes to not be zero")
		}
		want := flatTypes[0]
		flatTypes = flatTypes[1:]
		switch {
		case have.Kind() == kind.Float32 && want == kind.U32:
			f32, ok := have.Value().(float32)
			if !ok {
				return nil, types.NewCastError(have.Value(), "float32")
			}
			u32 := math.Float32bits(f32)
			payload[i] = values.U32(u32)
		case have.Kind() == kind.U32 && want == kind.U64:
			u32, ok := have.Value().(uint32)
			if !ok {
				return nil, types.NewCastError(have.Value(), "uint64")
			}
			payload[i] = values.U64(u32)
		case have.Kind() == kind.Float32 && want == kind.U64:
			f32, ok := have.Value().(float32)
			if !ok {
				return nil, types.NewCastError(have.Value(), "float32")
			}
			u32 := math.Float32bits(f32)
			payload[i] = values.U64(u32)
		case have.Kind() == kind.Float64 && want == kind.U64:
			f64, ok := have.Value().(float64)
			if !ok {
				return nil, types.NewCastError(have.Value(), "float64")
			}
			u64 := math.Float64bits(f64)
			payload[i] = values.U64(u64)
		default:
		}
	}
	for _, want := range flatTypes {
		zero, err := values.Zero(want)
		if err != nil {
			return nil, err
		}
		payload = append(payload, zero)
	}
	return append([]values.Value{values.U32(caseIndex)}, payload...), nil
}

func MatchCase(v any, cases []types.Case) (uint32, any, error) {
	vMap, ok := v.(map[string]any)
	if !ok {
		return 0, nil, types.NewCastError(v, "map[string]any")
	}
	if len(vMap) != 1 {
		return 0, nil, fmt.Errorf("expected map with one element")
	}
	var key string
	var value any
	for key, value = range vMap {
	}

	labelMap := map[string]int{}
	for i, c := range cases {
		labelMap[c.Label] = i
	}
	for _, label := range strings.Split(key, "|") {
		caseIndex, ok := labelMap[label]
		if ok {
			return uint32(caseIndex), value, nil
		}
	}
	return 0, nil, fmt.Errorf("unable to locate label in cases")
}

func slice(values ...values.Value) []values.Value {
	return values
}
