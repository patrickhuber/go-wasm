package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func LiftOwn(cx *types.Context, i uint32, t *types.Own) (*types.Handle, error) {
	return cx.Instance.Handles.Transfer(i, t)
}

func LiftBorrow(cx *types.Context, i uint32, t *types.Borrow) (*types.Handle, error) {
	h, err := cx.Instance.Handles.Get(i, t.ResourceType)
	if err != nil {
		return nil, err
	}
	h.LendCount += 1
	cx.Lenders = append(cx.Lenders, h)
	return &types.Handle{
		Rep:       h.Rep,
		LendCount: 0,
		Context:   nil,
	}, nil
}

func LiftFlat(cx *types.Context, vi values.ValueIterator, t types.ValType) (any, error) {
	t = t.Despecialize()
	k := t.Kind()
	switch k {
	case kind.Bool:
		return LiftFlatBool(vi)
	case kind.U8:
		return LiftFlatU8(vi)
	case kind.U16:
		return LiftFlatU16(vi)
	case kind.U32:
		return LiftFlatU32(vi)
	case kind.U64:
		return LiftFlatU64(vi)
	case kind.S8:
		return LiftFlatS8(vi)
	case kind.S16:
		return LiftFlatS16(vi)
	case kind.S32:
		return LiftFlatS32(vi)
	case kind.S64:
		return LiftFlatS64(vi)
	case kind.Float32:
		return LiftFlatFloat32(vi)
	case kind.Float64:
		return LiftFlatFloat64(vi)
	case kind.Char:
		return LiftFlatChar(vi)
	case kind.String:
		return LiftFlatString(cx, vi)
	case kind.List:
		return LiftFlatList(cx, vi, t)
	case kind.Record:
		r, ok := t.(*types.Record)
		if !ok {
			return nil, types.NewCastError(t, "*types.Record")
		}
		return LiftFlatRecord(cx, vi, r.Fields)
	case kind.Variant:
		v, ok := t.(*types.Variant)
		if !ok {
			return nil, types.NewCastError(t, "*types.Variant")
		}
		return LiftFlatVariant(cx, vi, v)
	case kind.Flags:
		f, ok := t.(*types.Flags)
		if !ok {
			return nil, types.NewCastError(t, "*types.Flags")
		}
		return LiftFlatFlags(vi, f)
	case kind.Own:
		v, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}
		o, ok := t.(*types.Own)
		if !ok {
			return nil, types.Trap()
		}
		i, ok := v.(uint32)
		if !ok {
			return nil, types.Trap()
		}
		return LiftOwn(cx, i, o)
	case kind.Borrow:

		v, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}

		b, ok := t.(*types.Borrow)
		if !ok {
			return nil, types.Trap()
		}

		i, ok := v.(uint32)
		if !ok {
			return nil, types.Trap()
		}
		return LiftBorrow(cx, i, b)
	}
	return nil, fmt.Errorf("unable to lift type %s", k)
}

func LiftFlatBool(vi values.ValueIterator) (bool, error) {
	b, err := vi.Next(kind.S32)
	if err != nil {
		return false, err
	}
	u32, ok := b.(uint32)
	if !ok {
		return false, types.NewCastError(b, "uint32")
	}
	return u32 == 1, nil
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
	s32, ok := i.(int32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return int8(s32), nil
}

func LiftFlatS16(vi values.ValueIterator) (int16, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	s32, ok := i.(int32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return int16(s32), nil
}

func LiftFlatS32(vi values.ValueIterator) (int32, error) {
	i, err := vi.Next(kind.U32)
	if err != nil {
		return 0, err
	}
	s32, ok := i.(int32)
	if !ok {
		return 0, types.NewCastError(i, "int32")
	}
	return s32, nil
}

func LiftFlatS64(vi values.ValueIterator) (int64, error) {
	i, err := vi.Next(kind.U64)
	if err != nil {
		return 0, err
	}
	s64, ok := i.(int64)
	if !ok {
		return 0, types.NewCastError(i, "int64")
	}
	return s64, nil
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
	// range check?
	return rune(u32), nil
}

func LiftFlatList(cx *types.Context, vi values.ValueIterator, t types.ValType) (any, error) {
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

func LiftFlatString(cx *types.Context, vi values.ValueIterator) (any, error) {
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

func LiftFlatRecord(cx *types.Context, vi values.ValueIterator, fields []types.Field) (any, error) {
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

func LiftFlatVariant(cx *types.Context, vi values.ValueIterator, variant *types.Variant) (any, error) {
	flatTypes, err := variant.Flatten()
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

	err = types.TrapIf(int(u32CaseIndex) >= len(variant.Cases))
	if err != nil {
		return nil, fmt.Errorf("expected case length to be less than %d %w", u32CaseIndex, err)
	}

	c := variant.Cases[u32CaseIndex]
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
		variant.CaseLabelWithRefinements(c): v,
	}, nil
}

func LiftFlatFlags(vi values.ValueIterator, f *types.Flags) (any, error) {
	flat := 0
	shift := 0
	numFlags := f.NumI32Flags()
	for i := 0; i < int(numFlags); i++ {
		next, err := vi.Next(kind.U32)
		if err != nil {
			return nil, err
		}
		u32Next, ok := next.(uint32)
		if !ok {
			return nil, types.NewCastError(next, "int32")
		}
		flat |= (int(u32Next) << shift)
		shift += 32
	}
	return UnpackFlagsFromInt(flat, f.Labels), nil
}
