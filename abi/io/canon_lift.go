package io

import (
	"fmt"

	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-types/tuple"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/trap"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/internal/collections"
)

func SetError[T any](res Result[T], err error) Result[T] {
	return NewError[T](err)
}

func CanonLift(opts *types.CanonicalOptions,
	inst *types.ComponentInstance,
	callee func(any) Result[any],
	ft types.FuncType,
	args []any,
	maxFlatParams int,
	maxFlatResults int) (res Result[Tuple2[[]any, types.PostReturnFunc]]) {

	// handle any try failures
	handle.Error(&res)

	if !inst.MayEnter {
		return SetError(res, types.TrapWith("ComponentInstance MayEnter must be true"))
	}
	if !inst.MayLeave {
		return SetError(res, fmt.Errorf("ComponentInstance MayLeave must be true"))
	}
	cx := &types.CallContext{
		Options:  opts,
		Instance: inst,
	}

	inst.MayLeave = false
	flatArgs := result.New(
		LowerValues(cx, maxFlatParams, args, ft.ParamTypes(), nil)).Unwrap()
	inst.MayLeave = true

	flatResults := callee(flatArgs).Unwrap()
	results, ok := flatResults.([]any)
	if !ok {
		SetError(res, types.NewCastError(flatResults, "[]any"))
	}

	valueResults := result.New(collections.Select[any, values.Value](results, func(source any) (values.Value, error) {
		v, ok := source.(values.Value)
		if !ok {
			return nil, types.NewCastError(source, "values.Value")
		}
		return v, nil
	})).Unwrap()

	lifted := LiftValues(cx, maxFlatResults, values.NewIterator(valueResults...), ft.ResultTypes()).Unwrap()

	var postResult types.PostReturnFunc = func() {
		if opts.PostReturn != nil {
			// python code has this set as:
			// opts.PostReturn(flatResults)
			opts.PostReturn()
		}
		cx.ExitCall()
	}
	return result.Ok(tuple.New2(lifted, postResult))
}

func LiftValues(cx *types.CallContext, maxFlat int, vi values.ValueIterator, ts []types.ValType) (res Result[[]any]) {
	defer handle.Error(&res)
	flatTypes := FlattenTypes(ts).Unwrap()
	if len(flatTypes) <= maxFlat {
		liftFlat := func(t types.ValType) (any, error) {
			return LiftFlat(cx, vi, t)
		}
		return collections.Select(ts, liftFlat)
	}

	ptr := Cast[any, uint32](vi.Next(kind.U32).Unwrap()).Unwrap()

	tupleType := types.NewTuple(ts...)

	alignment := Alignment(tupleType).Unwrap()
	aligned := AlignTo(ptr, alignment)

	trap.Iff(ptr != aligned, "ptr %d not aligned to %d", ptr, aligned)

	size := Size(tupleType).Unwrap()

	trap.Iff(ptr+size > uint32(cx.Options.Memory.Len()),
		"ptr %d + offset %d is greater than len(memory) %d",
		ptr, size, cx.Options.Memory.Len())

	load := Load(cx, tupleType, ptr).Unwrap()

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
