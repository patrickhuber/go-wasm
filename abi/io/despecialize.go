package io

import (
	"strconv"

	"github.com/patrickhuber/go-wasm/abi/types"
)

func Despecialize(t types.ValType) types.ValType {
	switch vt := t.(type) {
	case types.Tuple:
		fields := []types.Field{}
		for i, typ := range vt.Types() {
			fields = append(fields, types.Field{
				Label: strconv.Itoa(i),
				Type:  typ,
			})
		}
		return types.NewRecord(fields...)
	case types.Union:
		cases := []types.Case{}
		for i, typ := range vt.Types() {
			cases = append(cases, types.Case{Label: strconv.Itoa(i), Type: typ})
		}
		return types.NewVariant(cases...)
	case types.Enum:
		cases := []types.Case{}
		for _, label := range vt.Labels() {
			cases = append(cases, types.Case{Label: label, Type: nil})
		}
		return types.NewVariant(cases...)
	case types.Option:
		return types.NewVariant(
			types.Case{
				Label: "none",
				Type:  nil,
			},
			types.Case{
				Label: "some",
				Type:  vt,
			})

	case types.Result:
		return types.NewVariant(
			types.Case{
				Label: "ok",
				Type:  vt.Ok(),
			},
			types.Case{
				Label: "error",
				Type:  vt.Error(),
			})
	}
	return t
}
