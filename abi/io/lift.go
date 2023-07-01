package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func LiftOwn(cx *types.CallContext, i uint32, t types.Own) (*types.HandleElem, error) {
	return cx.Instance.Handles.Remove(i, t)
}

func LiftBorrow(cx *types.CallContext, i uint32, t types.Borrow) (*types.HandleElem, error) {
	h, err := cx.Instance.Handles.Get(i, t.ResourceType())
	if err != nil {
		return nil, err
	}
	h.LendCount += 1
	cx.Lenders = append(cx.Lenders, h)
	return &types.HandleElem{
		Rep:       h.Rep,
		LendCount: 0,
		Context:   nil,
	}, nil
}

func LiftFlat(cx *types.CallContext, vi values.ValueIterator, t types.ValType) (any, error) {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool:
		return LiftFlatBool(vi)
	case types.U8:
		return LiftFlatU8(vi)
	case types.U16:
		return LiftFlatU16(vi)
	case types.U32:
		return LiftFlatU32(vi)
	case types.U64:
		return LiftFlatU64(vi)
	case types.S8:
		return LiftFlatS8(vi)
	case types.S16:
		return LiftFlatS16(vi)
	case types.S32:
		return LiftFlatS32(vi)
	case types.S64:
		return LiftFlatS64(vi)
	case types.Float32:
		return LiftFlatFloat32(vi)
	case types.Float64:
		return LiftFlatFloat64(vi)
	case types.Char:
		return LiftFlatChar(vi)
	case types.String:
		return LiftFlatString(cx, vi)
	case types.List:
		return LiftFlatList(cx, vi, vt.Type())
	case types.Record:
		return LiftFlatRecord(cx, vi, vt.Fields())
	case types.Variant:
		return LiftFlatVariant(cx, vi, vt)
	case types.Flags:
		return LiftFlatFlags(vi, vt)
	case types.Own:
		v, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}
		i, ok := v.(uint32)
		if !ok {
			return nil, types.TrapWith("unable to cast %T to uint32", v)
		}
		return LiftOwn(cx, i, vt)
	case types.Borrow:

		v, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}

		i, ok := v.(uint32)
		if !ok {
			return nil, types.TrapWith("unable to cast %T to uint32", v)
		}
		return LiftBorrow(cx, i, vt)
	}
	return nil, fmt.Errorf("LiftFlat: unable to match type %T", t)
}

func LiftFlatBool(vi values.ValueIterator) (bool, error) {
	b, err := vi.Next(kind.U32)
	if err != nil {
		return false, err
	}
	u32, ok := b.(uint32)
	if !ok {
		return false, types.NewCastError(b, "uint32")
	}
	return u32 != 0, nil
}

func LiftFlatU8(vi values.ValueIterator) (uint8, error) {
	// s8 is packed as a s32
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "uint32")
	}
	return uint8(u32), nil
}

func LiftFlatU16(vi values.ValueIterator) (uint16, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "uint32")
	}
	return uint16(u32), nil
}

func LiftFlatU32(vi values.ValueIterator) (uint32, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "uint32")
	}
	return u32, nil
}

func LiftFlatU64(vi values.ValueIterator) (uint64, error) {
	i, err := vi.Next(kind.U64)
	if err != nil {
		return 0, err
	}
	u64, ok := i.(uint64)
	if !ok {
		return 0, types.NewCastError(i, "uint64")
	}
	return u64, nil
}

func LiftFlatS8(vi values.ValueIterator) (int8, error) {
	// s8 is packed as a s32
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return int8(u32), nil
}

func LiftFlatS16(vi values.ValueIterator) (int16, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return int16(u32), nil
}

func LiftFlatS32(vi values.ValueIterator) (int32, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	u32, ok := i.(uint32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return int32(u32), nil
}

func LiftFlatS64(vi values.ValueIterator) (int64, error) {
	i, err := vi.Next(kind.U64)
	if err != nil {
		return 0, err
	}
	u64, ok := i.(uint64)
	if !ok {
		return 0, types.NewCastError(i, "int64")
	}
	return int64(u64), nil
}

func LiftFlatFloat32(vi values.ValueIterator) (float32, error) {
	f, err := vi.Next(kind.Float32)
	if err != nil {
		return 0, err
	}
	f32, ok := f.(float32)
	if !ok {
		return 0, types.NewCastError(f, "float32")
	}
	return f32, nil
}

func LiftFlatFloat64(vi values.ValueIterator) (float64, error) {
	f, err := vi.Next(kind.Float64)
	if err != nil {
		return 0, err
	}
	f64, ok := f.(float64)
	if !ok {
		return 0, types.NewCastError(f, "float64")
	}
	return f64, nil
}

func LiftFlatChar(vi values.ValueIterator) (rune, error) {
	u32, err := LiftFlatU32(vi)
	var r rune
	if err != nil {
		return r, err
	}
	return ConvertU32ToRune(u32)
}

func LiftFlatList(cx *types.CallContext, vi values.ValueIterator, t types.ValType) (any, error) {
	ptr, err := LiftFlatU32(vi)
	if err != nil {
		return nil, err
	}
	length, err := LiftFlatU32(vi)
	if err != nil {
		return nil, err
	}
	return LoadListFromRange(cx, ptr, length, t)
}

func LiftFlatString(cx *types.CallContext, vi values.ValueIterator) (any, error) {
	ptr, err := LiftFlatU32(vi)
	if err != nil {
		return nil, err
	}
	packedLength, err := LiftFlatU32(vi)
	if err != nil {
		return nil, err
	}
	return LoadStringFromRange(cx, ptr, packedLength)
}

func LiftFlatRecord(cx *types.CallContext, vi values.ValueIterator, fields []types.Field) (any, error) {
	record := map[string]any{}
	for _, f := range fields {
		value, err := LiftFlat(cx, vi, f.Type)
		if err != nil {
			return nil, err
		}
		record[f.Label] = value
	}
	return record, nil
}

func LiftFlatVariant(cx *types.CallContext, vi values.ValueIterator, variant types.Variant) (any, error) {
	flatTypes, err := FlattenType(variant)
	if err != nil {
		return nil, err
	}
	if len(flatTypes) == 0 {
		return nil, fmt.Errorf("expected at least one type found 0")
	}

	first := flatTypes[0]
	flatTypes = flatTypes[1:]

	if first != kind.U32 {
		return nil, fmt.Errorf("expected kind.U32 found kind.%s", first.String())
	}

	caseIndex, err := vi.Next(kind.U32)
	if err != nil {
		return nil, err
	}

	u32CaseIndex, ok := caseIndex.(uint32)
	if !ok {
		return nil, types.NewCastError(caseIndex, "uint32")
	}

	if int(u32CaseIndex) >= len(variant.Cases()) {
		return nil, types.TrapWith("case index %d exceeds bounds of cases %d", u32CaseIndex, len(variant.Cases()))
	}

	c := variant.Cases()[u32CaseIndex]
	var v any
	if c.Type == nil {
		v = nil
	} else {
		cvi := values.NewCoerceValueIterator(vi, flatTypes)
		v, err = LiftFlat(cx, cvi, c.Type)
		if err != nil {
			return nil, err
		}
		flatTypes = cvi.FlatTypes()
		vi = cvi.ValueIterator()
	}

	for _, have := range flatTypes {
		_, err := vi.Next(have)
		if err != nil {
			return nil, err
		}
	}
	return map[string]any{
		CaseLabelWithRefinements(c, variant.Cases()): v,
	}, nil
}

func LiftFlatFlags(vi values.ValueIterator, f types.Flags) (any, error) {
	var flat uint64 = 0
	shift := 0
	numFlags := NumI32Flags(f.Labels())
	for i := 0; i < int(numFlags); i++ {
		next, err := vi.Next(kind.U32)
		if err != nil {
			return nil, err
		}
		u32Next, ok := next.(uint32)
		if !ok {
			return nil, types.NewCastError(next, "int32")
		}
		flat |= (uint64(u32Next) << shift)
		shift += 32
	}
	return UnpackFlagsFromInt(flat, f.Labels()), nil
}
