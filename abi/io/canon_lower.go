package io

import (
	"fmt"
	"strconv"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/abi/values"
	"github.com/patrickhuber/go-wasm/internal/collections"
)

func CanonLower(
	opts *types.CanonicalOptions,
	inst *types.ComponentInstance,
	callee func([]any) ([]any, types.PostReturnFunc, error),
	callingImport bool,
	ft types.FuncType,
	flatArgs []any,
	maxFlatParams int,
	maxFlatResults int) ([]any, error) {

	cx := &types.CallContext{
		Options:  opts,
		Instance: inst,
	}
	if !inst.MayLeave {
		return nil, types.TrapWith("CanonLower : ComponentInstance MayEnter must be true")
	}
	if !inst.MayEnter {
		return nil, fmt.Errorf("CanonLower : ComponentInstance MayLeave must be true")
	}
	args, err := collections.Select(flatArgs, func(source any) (values.Value, error) {
		result, ok := any(source).(values.Value)
		if !ok {
			return nil, types.NewCastError(source, "values.Value")
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	vi := values.NewIterator(args...)
	lifted, err := LiftValues(cx, maxFlatParams, vi, ft.ParamTypes())
	if err != nil {
		return nil, err
	}

	results, postReturn, err := callee(lifted)
	if err != nil {
		return nil, err
	}

	inst.MayLeave = false
	flatResults, err := LowerValues(cx, maxFlatResults, results, ft.ResultTypes(), vi)
	if err != nil {
		return nil, err
	}
	inst.MayLeave = true

	postReturn()
	cx.ExitCall()

	if callingImport {
		inst.MayEnter = true
	}
	return flatResults, nil
}

func LowerValues(cx *types.CallContext, maxFlat int, vs []any, ts []types.ValType, outParam values.ValueIterator) ([]any, error) {
	flatTypes, err := FlattenTypes(ts)
	if err != nil {
		return nil, err
	}

	if len(flatTypes) <= maxFlat {
		flatVals := []any{}
		for i, v := range vs {
			flat, err := LowerFlat(cx, v, ts[i])
			if err != nil {
				return nil, err
			}
			for _, v := range flat {
				flatVals = append(flatVals, v)
			}
		}
		return flatVals, nil
	}

	return LowerValuesToTuple(ts, vs, outParam, cx)
}

func LowerValuesToTuple(ts []types.ValType, vs []any, outParam values.ValueIterator, cx *types.CallContext) ([]any, error) {
	tupleType := types.NewTuple(ts...)
	tupleValue := map[string]any{}
	for i, v := range vs {
		tupleValue[strconv.Itoa(i)] = v
	}

	alignment, err := Alignment(tupleType)
	if err != nil {
		return nil, err
	}

	size, err := Size(tupleType)
	if err != nil {
		return nil, err
	}

	var ptr uint32
	if outParam == nil {
		ptr, err = cx.Options.Realloc(0, 0, alignment, size)
		if err != nil {
			return nil, err
		}

	} else {
		p, err := outParam.Next(kind.U32)
		if err != nil {
			return nil, err
		}
		var ok bool
		ptr, ok = p.(uint32)
		if !ok {
			return nil, types.NewCastError(p, "uint32")
		}
	}
	align, err := AlignTo(ptr, alignment)
	if err != nil {
		return nil, err
	}
	if ptr != align {
		return nil, fmt.Errorf("ptr %d does not align to %d", ptr, align)
	}
	if ptr+size > uint32(cx.Options.Memory.Len()) {
		return nil, fmt.Errorf("ptr %d is greater than memory size %d", ptr+size, cx.Options.Memory.Len())
	}
	return []any{ptr}, nil
}
