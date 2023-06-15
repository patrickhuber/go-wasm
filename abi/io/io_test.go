package io_test

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
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
		}, []any{uint32(1), uint32(2), uint32(3)}, map[string]any{"x": uint8(1), "y": uint16(2), "z": uint32(3)}, encoding.None, nil, nil},
		{"tuple", tuple(
			tuple(&types.U8{}, &types.U8{}),
			&types.U8{}), []any{uint32(1), uint32(2), uint32(3)}, map[string]any{"0": map[string]any{"0": uint8(1), "1": uint8(2)}, "1": uint8(3)}, encoding.UTF8, nil, nil},
		{"flags", flags(), []any{}, map[string]any{}, encoding.UTF8, nil, nil},
		{"flags", flags("a", "b"), []any{uint32(0)}, map[string]any{"a": false, "b": false}, encoding.UTF8, nil, nil},
		{"flags", flags("a", "b"), []any{uint32(2)}, map[string]any{"a": false, "b": true}, encoding.UTF8, nil, nil},
		{"flags", flags("a", "b"), []any{uint32(3)}, map[string]any{"a": true, "b": true}, encoding.UTF8, nil, nil},
		{"flags", flags("a", "b"), []any{uint32(4)}, map[string]any{"a": false, "b": false}, encoding.UTF8, nil, nil},
		{"flags", flags(Apply(Range(0, 33), strconv.Itoa)...), []any{uint32(math.MaxUint32), uint32(0x1)}, Zip(Apply(Range(0, 33), strconv.Itoa), Repeat[any](true, 33)), encoding.UTF8, nil, nil},
		{"variant", variant(vcase("x", &types.U8{}, nil), vcase("y", &types.Float32{}, nil), vcase("z", nil, nil)), []any{uint32(0), uint32(42)}, map[string]any{"x": uint8(42)}, encoding.UTF8, nil, nil},
		{"variant", variant(vcase("x", &types.U8{}, nil), vcase("y", &types.Float32{}, nil), vcase("z", nil, nil)), []any{uint32(0), uint32(256)}, map[string]any{"x": uint8(0)}, encoding.UTF8, nil, nil},
		{"variant", variant(vcase("x", &types.U8{}, nil), vcase("y", &types.Float32{}, nil), vcase("z", nil, nil)), []any{uint32(1), uint32(0x4048f5c3)}, map[string]any{"y": float32(3.140000104904175)}, encoding.UTF8, nil, nil},
		{"variant", variant(vcase("x", &types.U8{}, nil), vcase("y", &types.Float32{}, nil), vcase("z", nil, nil)), []any{uint32(2), uint32(0xffffffff)}, map[string]any{"z": nil}, encoding.UTF8, nil, nil},
	}

	for _, oneTest := range tests {
		t.Run(oneTest.name, func(t *testing.T) {
			cxt := NewContext(&bytes.Buffer{}, encoding.UTF8, nil, nil)
			err := test(oneTest.t, oneTest.valsToLift, oneTest.v, cxt, oneTest.dstEncoding, oneTest.lowerT, oneTest.lowerV)
			require.Nil(t, err)
		})
	}
}

func flags(labels ...string) *types.Flags {
	return &types.Flags{
		Labels: labels,
	}
}

func tuple(t ...types.ValType) *types.Tuple {
	return &types.Tuple{
		Types: t,
	}
}

func variant(c ...types.Case) *types.Variant {
	return &types.Variant{
		Cases: c,
	}
}

func vcase(label string, val types.ValType, refines *string) types.Case {
	return types.Case{
		Label:   label,
		Type:    val,
		Refines: refines,
	}
}

func Range(low int, high int) []int {
	var result []int
	for i := low; i < high; i++ {
		result = append(result, i)
	}
	return result
}

func Repeat[TValue any](value TValue, times int) []TValue {
	var result []TValue
	for i := 0; i < times; i++ {
		result = append(result, value)
	}
	return result
}

func Apply[TInput any, TOutput any](input []TInput, transform func(k TInput) TOutput) []TOutput {
	var slice []TOutput
	for _, i := range input {
		slice = append(slice, transform(i))
	}
	return slice
}

func Zip[TKey comparable, TValue any](keys []TKey, values []TValue) map[TKey]TValue {
	length := len(keys)
	if len(values) < length {
		length = len(values)
	}
	result := make(map[TKey]TValue)
	for i := 0; i < length; i++ {
		result[keys[i]] = values[i]
	}
	return result
}

func Slice[T any](values ...T) []T {
	return values
}

func test(t types.ValType, valsToLift []any, v any,
	cx *types.Context,
	dstEncoding encoding.Encoding, lowerT types.ValType, lowerV any) error {

	flattened, err := t.Flatten()
	if err != nil {
		return err
	}
	vs, err := zip(flattened, valsToLift)
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
		case kind.U32:
			vals, err = io.LowerU32(v[i])
		case kind.U64:
			vals, err = io.LowerU64(v[i])
		case kind.Float32:
			vals, err = io.LowerFloat32(v[i])
		case kind.Float64:
			vals, err = io.LowerFloat64(v[i])
		default:
			err = fmt.Errorf("invalid kind.%s", t.String())
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
