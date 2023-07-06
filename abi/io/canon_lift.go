package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/internal/collections"
)

func CanonLift(
	opts *types.CanonicalOptions,
	inst *types.ComponentInstance,
	callee func(any) (any, error),
	ft types.FuncType,
	args []any,
	maxFlatParams int,
	maxFlatResults int) ([]any, types.PostReturnFunc, error) {

	if !inst.MayEnter {
		return nil, nil, types.TrapWith("ComponentInstance MayEnter must be true")
	}
	if !inst.MayLeave {
		return nil, nil, fmt.Errorf("ComponentInstance MayLeave must be true")
	}
	cx := &types.CallContext{
		Options:  opts,
		Instance: inst,
	}

	inst.MayLeave = false

	flatArgs, err := LowerValues(cx, maxFlatParams, args, ft.ParamTypes(), nil)
	if err != nil {
		return nil, nil, err
	}

	inst.MayLeave = true

	flatResults, err := callee(flatArgs)
	if err != nil {
		return nil, nil, err
	}

	results, ok := flatResults.([]any)
	if !ok {
		return nil, nil, types.NewCastError(flatResults, "[]any")
	}

	valueResults, err := collections.Select[any, values.Value](results, func(source any) (values.Value, error) {
		v, ok := source.(values.Value)
		if !ok {
			return nil, types.NewCastError(source, "values.Value")
		}
		return v, nil
	})
	if err != nil {
		return nil, nil, err
	}

	lifted, err := LiftValues(cx, maxFlatResults, values.NewIterator(valueResults...), ft.ResultTypes())
	if err != nil {
		return nil, nil, err
	}

	postResult := func() {
		if opts.PostReturn != nil {
			// python code has this set as:
			// opts.PostReturn(flatResults)
			opts.PostReturn()
		}
		cx.ExitCall()
	}

	return lifted, postResult, nil
}

func LiftValues(cx *types.CallContext, maxFlat int, vi values.ValueIterator, ts []types.ValType) ([]any, error) {
	flatTypes, err := FlattenTypes(ts)

	if err != nil {
		return nil, err
	}

	if len(flatTypes) <= maxFlat {
		liftFlat := func(t types.ValType) (any, error) {
			return LiftFlat(cx, vi, t)
		}
		return collections.Select(ts, liftFlat)
	}

	anyPtr, err := vi.Next(kind.U32)
	if err != nil {
		return nil, err
	}

	ptr, ok := anyPtr.(uint32)
	if !ok {
		return nil, types.NewCastError(anyPtr, "uint32")
	}

	tupleType := types.NewTuple(ts...)

	alignment, err := Alignment(tupleType)
	if err != nil {
		return nil, err
	}

	aligned, err := AlignTo(ptr, alignment)
	if err != nil {
		return nil, err
	}

	if ptr != aligned {
		return nil, types.TrapWith("ptr %d not aligned to %d", ptr, aligned)
	}

	size, err := Size(tupleType)
	if err != nil {
		return nil, err
	}

	if ptr+size > uint32(cx.Options.Memory.Len()) {
		return nil, types.TrapWith("ptr %d + offset %d is greater than len(memory) %d", ptr, size, cx.Options.Memory.Len())
	}

	load, err := Load(cx, tupleType, ptr)
	if err != nil {
		return nil, err
	}

	m, ok := load.(map[string]any)
	if !ok {
		return nil, types.NewCastError(load, "map[string]any")
	}

	values := []any{}
	for _, v := range m {
		values = append(values, v)
	}
	return values, nil
}
