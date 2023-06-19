package wat_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/wat"
	"github.com/stretchr/testify/require"
)

func TestCanParse(t *testing.T) {
	type test struct {
		name   string
		wat    string
		module *wat.Module
	}
	tests := []test{
		{"module", "(module)", &wat.Module{}},
		{"memory", "(module (memory 1) (func))", wat.NewModule(func(s wat.SectionBuilder) {
			s.Function(func(f wat.FunctionBuilder) {})
			s.Memory(func(m wat.MemoryBuilder) {
				m.Limits(1)
			})
		}).Build()},
		{"func_alias", "(module (func $alias ))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.ID("$alias")
				})
			}).Build()},
		{"param_i32", "(module (func (param i32)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Parameters(func(p wat.ParametersBuilder) {
						p.Parameter(wat.I32)
					})
				})
			}).Build()},
		{"param_i64", "(module (func (param i64)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Parameters(func(p wat.ParametersBuilder) {
						p.Parameter(wat.I64)
					})
				})
			}).Build()},
		{"param_f32", "(module (func (param f32)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Parameters(func(p wat.ParametersBuilder) {
						p.Parameter(wat.F32)
					})
				})
			}).Build()},
		{"param_f64", "(module (func (param f64)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Parameters(func(p wat.ParametersBuilder) {
						p.Parameter(wat.F64)
					})
				})
			}).Build()},
		{"result_i32", "(module (func (result i32)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Results(func(p wat.ResultsBuilder) {
						p.Result(wat.I32)
					})
				})
			}).Build()},
		{"result_i64", "(module (func (result i64)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Results(func(p wat.ResultsBuilder) {
						p.Result(wat.I64)
					})
				})
			}).Build()},
		{"result_f32", "(module (func (result f32)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Results(func(p wat.ResultsBuilder) {
						p.Result(wat.F32)
					})
				})
			}).Build()},
		{"result_f64", "(module (func (result f64)))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Results(func(p wat.ResultsBuilder) {
						p.Result(wat.F64)
					})
				})
			}).Build()},
		{"multi_result", "(module (func (result i64) (result i64) ))", wat.NewModule(
			func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Results(func(p wat.ResultsBuilder) {
						p.Result(wat.I64)
						p.Result(wat.I64)
					})
				})
			}).Build()},
		{"function", `
			(module
				(func (param $lhs i32) (param $rhs i32) (result i32)
				  local.get $lhs
				  local.get $rhs
				  i32.add))`,
			wat.NewModule(func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {
					f.Parameters(func(p wat.ParametersBuilder) {
						p.Parameter(wat.I32).ID("$lhs")
						p.Parameter(wat.I32).ID("$rhs")
					})
					f.Results(func(r wat.ResultsBuilder) {
						r.Result(wat.I32)
					})
					f.Instructions(func(i wat.InstructionsBuilder) {
						i.Local(wat.LocalGet).ID("$lhs")
						i.Local(wat.LocalGet).ID("$rhs")
						i.I32Add()
					})
				})
			}).Build()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			CanParse(t, test.wat, test.module)
		})
	}
}

func CanParse(t *testing.T, input string, expected *wat.Module) {
	result, err := wat.ParseString(input)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, result, expected)
}
