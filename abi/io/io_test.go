package io_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	type testCase struct {
		name        string
		t           types.ValType
		valsToLift  []any
		v           any
		dstEncoding encoding.Encoding
		lowerT      types.ValType
		lowerV      any
	}
	tests := []testCase{
		{"record", &types.Record{}, []any{}, map[string]any{}, encoding.None, nil, nil},
		{"record_fields", &types.Record{
			Fields: []types.Field{
				{Label: "x", Type: &types.U8{}},
				{Label: "y", Type: &types.U16{}},
				{Label: "z", Type: &types.U32{}},
			},
		}, []any{int32(1), int32(2), int32(3)}, map[string]any{"x": uint8(1), "y": uint16(2), "z": uint32(3)}, encoding.None, nil, nil},
		{"tuple", &types.Tuple{
			Types: []types.ValType{
				&types.Tuple{
					Types: []types.ValType{
						&types.U8{},
						&types.U8{},
					},
				},
				&types.U8{},
			},
		}, []any{int32(1), int32(2), int32(3)}, map[string]any{"0": map[string]any{"0": uint8(1), "1": uint8(2)}, "1": uint8(3)}, encoding.UTF8, nil, nil},
		{"flags", &types.Flags{}, []any{}, map[string]any{}, encoding.UTF8, nil, nil},
		{"flags", &types.Flags{Labels: []string{"a", "b"}}, []any{int32(0)}, map[string]any{"a": false, "b": false}, encoding.UTF8, nil, nil},
		{"flags", &types.Flags{Labels: []string{"a", "b"}}, []any{int32(2)}, map[string]any{"a": false, "b": true}, encoding.UTF8, nil, nil},
		{"flags", &types.Flags{Labels: []string{"a", "b"}}, []any{int32(3)}, map[string]any{"a": true, "b": true}, encoding.UTF8, nil, nil},
		{"flags", &types.Flags{Labels: []string{"a", "b"}}, []any{int32(4)}, map[string]any{"a": false, "b": false}, encoding.UTF8, nil, nil},
	}
	for _, oneTest := range tests {
		t.Run(oneTest.name, func(t *testing.T) {
			cxt := NewContext(&bytes.Buffer{}, encoding.UTF8, nil, nil)
			err := test(oneTest.t, oneTest.valsToLift, oneTest.v, cxt, oneTest.dstEncoding, oneTest.lowerT, oneTest.lowerV)
			require.Nil(t, err)
		})
	}
}

func test(t types.ValType, valsToLift []any, v any,
	cx *types.Context,
	dstEncoding encoding.Encoding, lowerT types.ValType, lowerV any) error {

	vs, err := zip(t.Flatten(), valsToLift)
	if err != nil {
		return err
	}

	vi := values.NewIterator(vs...)

	got, err := io.LiftFlat(cx, vi, t)
	if err != nil {
		return fmt.Errorf("test : %w", err)
	}
	if v == nil {
		return fmt.Errorf("expected trap but got %v", got)
	}

	err = types.TrapIf(vi.Index() != vi.Length())
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(got, v) {
		return fmt.Errorf("initial lift_flat() expected %v but got %v", v, got)
	}

	lowerT = coalesce(lowerT, t)
	lowerV = coalesce(lowerV, v)

	heap := NewHeap(5 * cx.Options.Memory.Len())
	if dstEncoding == encoding.None {
		dstEncoding = cx.Options.StringEncoding
	}

	cx = NewContext(heap.Memory, dstEncoding, heap.ReAllocate, cx.Options.PostReturn)

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

	return reflect.DeepEqual(s, v)
}

func coalesce[T comparable](v T, other ...T) T {
	var zero T
	if v != zero {
		return v
	}
	for _, o := range other {
		if o != zero {
			return o
		}
	}
	return zero
}
