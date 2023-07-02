package io_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/stretchr/testify/require"
)

func TestFlatten(t *testing.T) {
	type test struct {
		name    string
		ft      types.FuncType
		params  []kind.Kind
		results []kind.Kind
	}
	params := []types.ValType{U8(), Float32(), Float64()}
	paramKinds := []kind.Kind{kind.U32, kind.Float32, kind.Float64}
	tests := []test{
		{"p8_pf32_pf64", FuncType(params, []types.ValType{}), paramKinds, []kind.Kind{}},
		{"p8_pf32_pf64_rf32", FuncType(params, []types.ValType{Float32()}), paramKinds, []kind.Kind{kind.Float32}},
		{"p8_pf32_pf64_ru8", FuncType(params, []types.ValType{U8()}), paramKinds, []kind.Kind{kind.U32}},
		{"p8_pf32_pf64_rtup_f32", FuncType(params, []types.ValType{Tuple(Float32())}), paramKinds, []kind.Kind{kind.Float32}},
		{"p8_pf32_pf64_rtup_f32_f32", FuncType(params, []types.ValType{Tuple(Float32(), Float32())}), paramKinds, []kind.Kind{kind.Float32, kind.Float32}},
		{"p8_pf32_pf64_rf32_rf32", FuncType(params, []types.ValType{Float32(), Float32()}), paramKinds, []kind.Kind{kind.Float32, kind.Float32}},

		{"pu8x17", FuncType(Repeat[types.ValType](U8(), 17), []types.ValType{}), Repeat(kind.U32, 17), []kind.Kind{}},
		{"pu8x17_rtup_u8_u8", FuncType(Repeat[types.ValType](U8(), 17), []types.ValType{Tuple(U8(), U8())}), Repeat(kind.U32, 17), Repeat(kind.U32, 2)},
	}
	const (
		MaxFlatResults = 1
		MaxFlatParams  = 16
	)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			params := test.params
			results := test.results

			if len(test.params) > MaxFlatParams {
				params = []kind.Kind{kind.U32}
			}
			if len(test.results) > MaxFlatResults {
				params = []kind.Kind{kind.U32}
			}
			expect := CoreFuncType(params, results)

			got, err := io.FlattenFuncTypeLift(test.ft, MaxFlatParams, MaxFlatResults)
			require.Nil(t, err)
			require.Equal(t, expect, got)

			if len(test.results) > MaxFlatResults {
				expect = CoreFuncType(append(test.params, kind.U32), []kind.Kind{})
			}

			got, err = io.FlattenFuncTypeLower(test.ft, MaxFlatParams, MaxFlatResults)
			require.Nil(t, err)
			require.Equal(t, expect, got)
		})
	}
}

func CoreFuncType(params []kind.Kind, results []kind.Kind) types.CoreFuncType {
	return types.NewCoreFuncType(params, results)
}
