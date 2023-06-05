package io_test

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
)

func test(t types.ValType, valsToLift []any, v any,
	cx *types.Context,
	dstEncoding types.StringEncoding, lowerT types.ValType, lowerV any) error {

	vs, err := zip(t.Flatten(), valsToLift)
	if err != nil {
		return err
	}
	vi := values.NewIterator(vs...)
	if v == nil {
		got, err := io.LiftFlat(cx, vi, t)
		if err != nil {
			return err
		}
		return fmt.Errorf("expected trap but got %v", got)
	}
	got, err := io.LiftFlat(cx, vi, t)
	if err != nil {
		return err
	}

	err = types.TrapIf(vi.Index() != vi.Length())
	if err != nil {
		return err
	}
	if got != v {
		return fmt.Errorf("initial lift_flat() expected %v but got %v", v, got)
	}

	if lowerT == nil {
		lowerT = t
	}
	if lowerV == nil {
		lowerV = v
	}

	heap := NewHeap(5 * cx.Options.Memory.Len())
	if dstEncoding == types.None {
		dstEncoding = cx.Options.StringEncoding
	}

	cx = NewContext(heap.Memory, types.StringEncoding(dstEncoding), heap.ReAllocate, cx.Options.PostReturn)

	loweredValues, err := io.LowerFlat(cx, v, lowerT)
	if err != nil {
		return err
	}
	// assert here with lowerT

	vi = values.NewIterator(loweredValues...)
	got, err = io.LiftFlat(cx, vi, lowerT)
	if err != nil {
		return err
	}
	if !equalModuloStringEncoding(got, lowerV) {
		return fmt.Errorf("re-lift expected %v but got %v", lowerV, got)
	}
	return nil
}

// zip emulates python zip but only used here
func zip(types []kind.Kind, v []any) ([]values.Value, error) {
	vs := []values.Value{}

	if len(v) != len(types) {
		return nil, fmt.Errorf("expected len(values)=%d to equal len(types)=%d", len(v), len(types))
	}
	for i := 0; i < len(v); i++ {
		t := types[i]
		var vals []values.Value
		var err error
		switch t {
		case kind.S32:
			vals, err = io.LowerS32(v[i])
		case kind.S64:
			vals, err = io.LowerS64(v[i])
		case kind.Float32:
			vals, err = io.LowerFloat32(v[i])
		case kind.Float64:
			vals, err = io.LowerFloat64(v[i])
		default:
			err = fmt.Errorf("invalid kind %s", t.String())
		}
		if err != nil {
			return nil, err
		}
		vs = append(vs, vals...)
	}
	return vs, nil
}

func equalModuloStringEncoding(s any, v any) bool {
	if s == nil && v == nil {
		return true
	}
	if s == nil {
		return false
	}
	if v == nil {
		return false
	}
	// TODO: list check
	// TODO: map check
	return s == v
}