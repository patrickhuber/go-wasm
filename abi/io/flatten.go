package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

func FlattenTypes(ts []types.ValType) ([]kind.Kind, error) {
	flat := []kind.Kind{}
	for _, t := range ts {
		flattened, err := FlattenType(t)
		if err != nil {
			return nil, err
		}
		flat = append(flat, flattened...)
	}
	return flat, nil
}

func FlattenType(t types.ValType) ([]kind.Kind, error) {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool, types.U8, types.U16, types.U32:
		return []kind.Kind{kind.U32}, nil
	case types.U64, types.S64:
		return []kind.Kind{kind.U64}, nil
	case types.S8, types.S16, types.S32:
		return []kind.Kind{kind.U32}, nil
	case types.F32:
		return []kind.Kind{kind.Float32}, nil
	case types.F64:
		return []kind.Kind{kind.Float64}, nil
	case types.Char:
		return []kind.Kind{kind.U32}, nil
	case types.String:
		return []kind.Kind{kind.U32, kind.U32}, nil
	case types.List:
		return []kind.Kind{kind.U32, kind.U32}, nil
	case types.Record:
		return FlattenRecord(vt)
	case types.Variant:
		return FlattenVariant(vt)
	case types.Flags:
		flat := []kind.Kind{}
		n := NumI32Flags(vt.Labels())
		for i := uint32(0); i < n; i++ {
			flat = append(flat, kind.U32)
		}
		return flat, nil
	case types.Own, types.Borrow:
		return []kind.Kind{kind.U32}, nil
	}
	return nil, fmt.Errorf("flatten_type: unable to match type %T", t)
}

func FlattenRecord(r types.Record) ([]kind.Kind, error) {
	flat := []kind.Kind{}
	for _, f := range r.Fields() {
		flattened, err := FlattenType(f.Type)
		if err != nil {
			return nil, err
		}
		flat = append(flat, flattened...)
	}
	return flat, nil
}

func FlattenVariant(v types.Variant) ([]kind.Kind, error) {
	flat := []kind.Kind{}
	for _, c := range v.Cases() {
		if c.Type == nil {
			continue
		}
		flattened, err := FlattenType(c.Type)
		if err != nil {
			return nil, err
		}
		for i, ft := range flattened {
			if i < len(flat) {
				flat[i] = join(flat[i], ft)
			} else {
				flat = append(flat, ft)
			}
		}
	}
	dt, err := DiscriminantType(v.Cases())
	if err != nil {
		return nil, err
	}

	flattened, err := FlattenType(dt)
	if err != nil {
		return nil, err
	}
	return append(flattened, flat...), nil
}

func FlattenFuncTypeLower(ft types.FuncType, maxFlatParams int, maxFlatResults int) (types.CoreFuncType, error) {
	flatParams, flatResults, err := flattenFuncType(ft, maxFlatParams)
	if err != nil {
		return nil, err
	}
	if len(flatResults) > maxFlatResults {
		flatParams = append(flatParams, kind.U32)
		flatResults = []kind.Kind{}
	}
	return types.NewCoreFuncType(flatParams, flatResults), nil
}

func FlattenFuncTypeLift(ft types.FuncType, maxFlatParams int, maxFlatResults int) (types.CoreFuncType, error) {
	flatParams, flatResults, err := flattenFuncType(ft, maxFlatParams)
	if err != nil {
		return nil, err
	}
	if len(flatResults) > maxFlatResults {
		flatResults = []kind.Kind{kind.U32}
	}
	return types.NewCoreFuncType(flatParams, flatResults), nil
}

func flattenFuncType(ft types.FuncType, maxFlatParams int) ([]kind.Kind, []kind.Kind, error) {
	flatParams, err := FlattenTypes(ft.ParamTypes())
	if err != nil {
		return nil, nil, err
	}
	if len(flatParams) > maxFlatParams {
		flatParams = []kind.Kind{kind.U32}
	}
	flatResults, err := FlattenTypes(ft.ResultTypes())
	if err != nil {
		return nil, nil, err
	}
	return flatParams, flatResults, nil
}

func join(a kind.Kind, b kind.Kind) kind.Kind {
	if a == b {
		return a
	}
	switch {
	case a == kind.U32 && b == kind.Float32:
		return kind.U32
	case a == kind.Float32 && b == kind.U32:
		return kind.U32
	default:
		return kind.U64
	}
}
