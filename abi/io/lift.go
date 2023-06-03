package io

import (
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
	switch t.Despecialize().Kind() {
	case kind.Bool:
		b, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}
		i32 := b.(int32)
		if i32 == 0 {
			return false, nil
		}
		return true, nil
	case kind.U8:
		return LiftFlatUnsigned(vi, 32, 8)
	case kind.U16:
		return LiftFlatUnsigned(vi, 32, 16)
	case kind.U32:
		return LiftFlatUnsigned(vi, 32, 32)
	case kind.U64:
		return LiftFlatUnsigned(vi, 64, 64)
	case kind.S8:
		return LiftFlatSigned(vi, 32, 8)
	case kind.S16:
		return LiftFlatSigned(vi, 32, 16)
	case kind.S32:
		return LiftFlatSigned(vi, 32, 32)
	case kind.S64:
		return LiftFlatSigned(vi, 64, 64)
	case kind.Float32:
		v, err := vi.Next(kind.Float32)
		if err != nil {
			return nil, err
		}
		return Canonicalize32(v)
	case kind.Float64:
		v, err := vi.Next(kind.Float64)
		if err != nil {
			return nil, err
		}
		return Canonicalize64(v)
	case kind.Char:
		v, err := vi.Next(kind.S32)
		if err != nil {
			return nil, err
		}
		return ConvertI32ToChar(cx, v)
	case kind.String:
		return LiftFlatString(cx, vi)
	case kind.List:
		return LiftFlatList(cx, vi, t)
	case kind.Record:
		r, ok := t.(*types.Record)
		if !ok {
			return nil, types.Trap()
		}
		return LiftFlatRecord(cx, vi, r.Fields)
	case kind.Variant:
		v, ok := t.(*types.Variant)
		if !ok {
			return nil, types.Trap()
		}
		return LiftFlatVariant(cx, vi, v.Cases)
	case kind.Flags:
		f, ok := t.(*types.Flags)
		if !ok {
			return nil, types.Trap()
		}
		return LiftFlatFlags(vi, f.Labels)
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
	return nil, types.Trap()
}

func LiftFlatList(cx *types.Context, vi values.ValueIterator, t types.ValType) (any, error) {
	panic("unimplemented")
}

func LiftFlatString(cx *types.Context, vi values.ValueIterator) (any, error) {
	panic("unimplemented")
}

func ConvertI32ToChar(cx *types.Context, any any) (any, error) {
	panic("unimplemented")
}

func Canonicalize64(any any) (any, error) {
	panic("unimplemented")
}

func Canonicalize32(any any) (any, error) {
	panic("unimplemented")
}

func LiftFlatSigned(vi values.ValueIterator, i1, i2 int) (any, error) {
	panic("unimplemented")
}

func LiftFlatUnsigned(vi values.ValueIterator, i1, i2 int) (any, error) {
	panic("unimplemented")
}
func LiftFlatRecord(cx *types.Context, vi values.ValueIterator, fields []types.Field) (any, error) {
	panic("unimplemented")
}
func LiftFlatVariant(cx *types.Context, vi values.ValueIterator, cases []types.Case) (any, error) {
	panic("unimplemented")
}
func LiftFlatFlags(vi values.ValueIterator, labels []string) (any, error) {
	panic("unimplemented")
}
