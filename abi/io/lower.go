package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func LowerFlat(cx *types.Context, v any, t types.ValType) ([]values.Value, error) {
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
	case kind.Record:
		r := t.(*types.Record)
		return LowerRecord(cx, v, r)
	}
	return nil, fmt.Errorf("unable to lower type %s", k.String())
}

func LowerBool(v any) ([]values.Value, error) {
	var i values.S32
	b, ok := v.(bool)
	if !ok {
		return nil, NewCastError(v, "bool")
	}
	if b {
		i = 1
	} else {
		i = 0
	}
	return slice(i), nil
}

func LowerU8(v any) ([]values.Value, error) {
	u8 := v.(uint8)
	i := values.S32(u8)
	return slice(i), nil
}

func LowerU16(v any) ([]values.Value, error) {
	u16 := v.(uint16)
	i := values.S32(u16)
	return slice(i), nil
}

func LowerU32(v any) ([]values.Value, error) {
	u32 := v.(uint32)
	i := values.S32(u32)
	return slice(i), nil
}

func LowerU64(v any) ([]values.Value, error) {
	u64 := v.(uint64)
	i := values.S64(u64)
	return slice(i), nil
}

func LowerS8(v any) ([]values.Value, error) {
	s8 := v.(int8)
	i := values.S32(s8)
	return slice(i), nil
}

func LowerS16(v any) ([]values.Value, error) {
	s16, ok := v.(int16)
	if !ok {
		return nil, NewCastError(v, "int16")
	}
	i := values.S32(s16)
	return slice(i), nil
}

func LowerS32(v any) ([]values.Value, error) {
	s32, ok := v.(int32)
	if !ok {
		return nil, NewCastError(v, "int32")
	}
	i := values.S32(s32)
	return slice(i), nil
}

func LowerS64(v any) ([]values.Value, error) {
	s64 := v.(int64)
	i := values.S64(s64)
	return slice(i), nil
}

func LowerFloat32(v any) ([]values.Value, error) {
	f32 := v.(float32)
	f := values.Float32(f32)
	return slice(f), nil
}

func LowerFloat64(v any) ([]values.Value, error) {
	f64 := v.(float64)
	f := values.Float64(f64)
	return slice(f), nil
}

func LowerChar(v any) ([]values.Value, error) {
	r := v.(rune)
	i := values.S32(r)
	return slice(i), nil
}

func LowerString(cx *types.Context, v any) ([]values.Value, error) {
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

func LowerRecord(cx *types.Context, v any, r *types.Record) ([]values.Value, error) {
	var flat []values.Value
	vMap, ok := v.(map[string]any)
	if !ok {
		return nil, NewCastError(v, "map[string]any")
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

func slice(values ...values.Value) []values.Value {
	return values
}
