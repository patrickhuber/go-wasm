package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/stretchr/testify/require"
)

func TestHandles(t *testing.T) {
	const (
		MaxFlatResults = 16
		MaxFlatParams  = 16
	)
	var dtorValue int
	dtor := func(x int) {
		dtorValue = x
	}
	inst := &types.ComponentInstance{}
	rt := types.NewResourceType(dtor, &types.ComponentInstance{})
	rt2 := types.NewResourceType(dtor, inst)
	opts := Options()
	hostImport := func(args []any) ([]any, types.PostReturnFunc, error) {
		require.Equal(t, 2, len(args), "args")
		require.Equal(t, 42, args[0])
		require.Equal(t, 44, args[1])
		return []any{uint32(45)}, func() {}, nil
	}
	coreWasm := func(val any) (any, error) {
		args, ok := val.([]any)
		if !ok {
			return nil, fmt.Errorf("args must be an array")
		}
		require.Equal(t, 4, len(args))

		var vs []values.Value
		for _, a := range args {
			v, ok := a.(values.Value)
			if !ok {
				return nil, types.NewCastError(a, "values.Value")
			}
			vs = append(vs, v)
		}
		expectedValues := []values.Value{
			values.U32(0),
			values.U32(1),
			values.U32(2),
			values.U32(3),
		}
		for i, v := range vs {
			require.Equal(t, expectedValues[i].Kind(), v.Kind())
			require.Equal(t, expectedValues[i].Value(), v.Value())
		}

		expectedResourceRepresentations := []uint32{42, 43, 44}
		for i, expected := range expectedResourceRepresentations {
			rep, err := io.CanonResourceRep(inst, rt, uint32(i))
			if err != nil {
				return nil, err
			}
			require.Equal(t, expected, rep)
		}

		hostFunctionType := FuncType(
			[]types.ValType{
				Borrow(rt),
				Borrow(rt),
			},
			[]types.ValType{
				Own(rt),
			})
		args = []any{
			values.U32(0),
			values.U32(2),
		}
		results, err := io.CanonLower(opts, inst, hostImport, true, hostFunctionType, args, MaxFlatParams, MaxFlatResults)
		if err != nil {
			return nil, err
		}
		require.Equal(t, 1, len(results))
		result0, ok := results[0].(values.Value)
		require.True(t, ok)
		require.Equal(t, kind.U32, result0.Kind())
		require.Equal(t, 3, result0.Value())
		rep, err := io.CanonResourceRep(inst, rt, 3)
		require.Nil(t, err)
		require.Equal(t, rep, 45)

		dtorValue = 0
		io.CanonResourceDrop(inst, rt, 0)

		return nil, nil
	}

	ft := FuncType(
		[]types.ValType{
			Own(rt),
			Own(rt),
			Borrow(rt),
			Borrow(rt2)},
		[]types.ValType{
			Own(rt),
			Own(rt),
			Own(rt)})
	args := []any{42, 43, 44, 13}
	got, _, err := io.CanonLift(opts, inst, coreWasm, ft, args, MaxFlatParams, MaxFlatResults)
	require.Nil(t, err)

	require.Equal(t, 3, len(got))
	require.Equal(t, 46, got[0])
	require.Equal(t, 43, got[1])
	require.Equal(t, 45, got[2])
	require.Equal(t, 4, len(inst.Handles.Table(rt).Array))
	for _, handle := range inst.Handles.Table(rt).Array {
		require.Nil(t, handle)
	}
	require.Equal(t, 4, len(inst.Handles.Table(rt).Free))
}
