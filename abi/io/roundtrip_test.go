package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/internal/collections"
	"github.com/stretchr/testify/require"
)

func TestCanRoundTrip(t *testing.T) {
	type test struct {
		name string
		t    types.ValType
		v    any
	}
	tests := []test{
		{"u8", S8(), int8(-1)},
		{"tuple_u16_u16", Tuple(U16(), U16()), NewTuple(uint16(3), uint16(4))},
		{"list_string", List(String()), []any{"hello there"}},
		{"list_list_string", List(List(String())), []any{[]any(nil), []any(nil)}},
		{"list_option_tuple_string_u16", List(Option(Tuple(String(), U16()))), []any{map[string]any{"some": NewTuple("answer", uint16(42))}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RoundTripTest(t, test.t, test.v)
		})
	}
}

func RoundTripTest(t *testing.T, typ types.ValType, v any) {
	const MaxFlatResults = 16
	const MaxFlatParams = 16
	ft := FuncType([]types.ValType{typ}, []types.ValType{typ})
	callee := func(val any) (any, error) { return val, nil }

	calleeHeap := NewHeap(1000)
	calleeOpts := Options(Memory(calleeHeap.Memory), Realloc(calleeHeap.ReAllocate))
	calleeInst := &types.ComponentInstance{MayEnter: true, MayLeave: true}
	liftedCallee := func(args []any) ([]any, types.PostReturnFunc, error) {
		return io.CanonLift(calleeOpts, calleeInst, callee, ft, args, MaxFlatParams, MaxFlatResults)
	}

	callerHeap := NewHeap(1000)
	callerOpts := Options(Memory(callerHeap.Memory), Realloc(callerHeap.ReAllocate))
	callerInst := &types.ComponentInstance{MayEnter: true, MayLeave: true}
	callerContext := &types.CallContext{Options: callerOpts, Instance: callerInst}

	flatArgs, err := io.LowerFlat(callerContext, v, typ)
	require.Nil(t, err)
	args := Select(flatArgs, func(vt values.Value) any { return vt })

	flatResults, err := io.CanonLower(callerOpts, callerInst, liftedCallee, true, ft, args, MaxFlatParams, MaxFlatResults)
	require.Nil(t, err)

	results, err := collections.Select(flatResults, Cast[any, values.Value])
	require.Nil(t, err)

	got, err := io.LiftFlat(callerContext, values.NewIterator(results...), typ)
	require.Nil(t, err)

	require.Equal(t, v, got)
	require.True(t, callerInst.MayLeave && callerInst.MayEnter)
	require.True(t, calleeInst.MayLeave && calleeInst.MayEnter)
}

func Cast[TSource, TTarget any](source TSource) (TTarget, error) {
	target, ok := any(source).(TTarget)
	if !ok {
		var zero TTarget
		return zero, types.NewCastError(source, fmt.Sprintf("%T", zero))
	}
	return target, nil
}

func Select[TSource, TTarget any](source []TSource, transform func(TSource) TTarget) []TTarget {
	target := []TTarget{}
	for _, s := range source {
		target = append(target, transform(s))
	}
	return target
}
