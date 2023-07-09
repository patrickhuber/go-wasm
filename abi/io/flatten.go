package io

import (
	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-types/tuple"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

func FlattenTypes(ts []types.ValType) (res Result[[]kind.Kind]) {
	defer handle.Error(&res)
	flat := []kind.Kind{}
	for _, t := range ts {
		flattened := FlattenType(t).Unwrap()
		flat = append(flat, flattened...)
	}
	return result.Ok(flat)
}

func FlattenType(t types.ValType) Result[[]kind.Kind] {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool, types.U8, types.U16, types.U32:
		return result.Ok([]kind.Kind{kind.U32})
	case types.U64, types.S64:
		return result.Ok([]kind.Kind{kind.U64})
	case types.S8, types.S16, types.S32:
		return result.Ok([]kind.Kind{kind.U32})
	case types.Float32:
		return result.Ok([]kind.Kind{kind.Float32})
	case types.Float64:
		return result.Ok([]kind.Kind{kind.Float64})
	case types.Char:
		return result.Ok([]kind.Kind{kind.U32})
	case types.String:
		return result.Ok([]kind.Kind{kind.U32})
	case types.List:
		return result.Ok([]kind.Kind{kind.U32})
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
		return result.Ok(flat)
	case types.Own, types.Borrow:
		return result.Ok([]kind.Kind{kind.U32})
	}
	return result.Errorf[[]kind.Kind]("flatten_type: unable to match type %T", t)
}

func FlattenRecord(r types.Record) Result[[]kind.Kind] {
	flat := []kind.Kind{}
	for _, f := range r.Fields() {
		flattened := FlattenType(f.Type).Unwrap()
		flat = append(flat, flattened...)
	}
	return result.Ok(flat)
}

func FlattenVariant(v types.Variant) (res Result[[]kind.Kind]) {
	defer handle.Error(&res)
	flat := []kind.Kind{}
	for _, c := range v.Cases() {
		if c.Type == nil {
			continue
		}
		flattened := FlattenType(c.Type).Unwrap()
		for i, ft := range flattened {
			if i < len(flat) {
				flat[i] = join(flat[i], ft)
			} else {
				flat = append(flat, ft)
			}
		}
	}
	dt := DiscriminantType(v.Cases()).Unwrap()
	flattened := FlattenType(dt).Unwrap()
	return result.Ok(append(flattened, flat...))
}

func FlattenFuncTypeLower(ft types.FuncType, maxFlatParams int, maxFlatResults int) (res Result[types.CoreFuncType]) {
	defer handle.Error(&res)
	flatParams, flatResults := flattenFuncType(ft, maxFlatParams).Unwrap().Deconstruct()

	if len(flatResults) > maxFlatResults {
		flatParams = append(flatParams, kind.U32)
		flatResults = []kind.Kind{}
	}
	return result.Ok(types.NewCoreFuncType(flatParams, flatResults))
}

func FlattenFuncTypeLift(ft types.FuncType, maxFlatParams int, maxFlatResults int) (res Result[types.CoreFuncType]) {
	defer handle.Error(&res)
	flat := flattenFuncType(ft, maxFlatParams).Unwrap()
	flatParams, flatResults := flat.Deconstruct()
	if len(flatResults) > maxFlatResults {
		flatResults = []kind.Kind{kind.U32}
	}
	return result.Ok(types.NewCoreFuncType(flatParams, flatResults))
}

func flattenFuncType(ft types.FuncType, maxFlatParams int) (res Result[Tuple2[[]kind.Kind, []kind.Kind]]) {
	defer handle.Error(&res)
	flatParams := FlattenTypes(ft.ParamTypes()).Unwrap()
	if len(flatParams) > maxFlatParams {
		flatParams = []kind.Kind{kind.U32}
	}
	flatResults := FlattenTypes(ft.ResultTypes()).Unwrap()
	return result.Ok(tuple.New2(flatParams, flatResults))
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
