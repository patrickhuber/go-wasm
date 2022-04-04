package wat_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wat"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		It("module", func() {
			canParse("(module)", &wat.Module{})
		})
		It("memory", func() {
			canParse("(module (memory 1) (func))", wat.NewModule(func(s wat.SectionBuilder) {
				s.Function(func(f wat.FunctionBuilder) {})
				s.Memory(func(m wat.MemoryBuilder) {
					m.Limits(1)
				})
			}).Build())
		})
		Describe("function", func() {
			It("can parse function alias", func() {
				canParse("(module (func $alias ))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.ID("$alias")
						})
					}).Build())
			})
			It("can parse i32 parameter", func() {
				canParse("(module (func (param i32)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Parameters(func(p wat.ParametersBuilder) {
								p.Parameter(wat.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 parameter", func() {
				canParse("(module (func (param i64)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Parameters(func(p wat.ParametersBuilder) {
								p.Parameter(wat.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 parameter", func() {
				canParse("(module (func (param f32)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Parameters(func(p wat.ParametersBuilder) {
								p.Parameter(wat.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 parameter", func() {
				canParse("(module (func (param f64)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Parameters(func(p wat.ParametersBuilder) {
								p.Parameter(wat.F64)
							})
						})
					}).Build())
			})
			It("can parse i32 result", func() {
				canParse("(module (func (result i32)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Results(func(p wat.ResultsBuilder) {
								p.Result(wat.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 result", func() {
				canParse("(module (func (result i64)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Results(func(p wat.ResultsBuilder) {
								p.Result(wat.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 result", func() {
				canParse("(module (func (result f32)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Results(func(p wat.ResultsBuilder) {
								p.Result(wat.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 result", func() {
				canParse("(module (func (result f64)))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Results(func(p wat.ResultsBuilder) {
								p.Result(wat.F64)
							})
						})
					}).Build())
			})
			It("can parse mutiple results", func() {
				canParse("(module (func (result i64) (result i64) ))", wat.NewModule(
					func(s wat.SectionBuilder) {
						s.Function(func(f wat.FunctionBuilder) {
							f.Results(func(p wat.ResultsBuilder) {
								p.Result(wat.I64)
								p.Result(wat.I64)
							})
						})
					}).Build())
			})
		})
		It("function", func() {
			builder := wat.NewModule(func(s wat.SectionBuilder) {
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
			})

			canParse(`
			(module
				(func (param $lhs i32) (param $rhs i32) (result i32)
				  local.get $lhs
				  local.get $rhs
				  i32.add))`,
				builder.Build())
		})
	})
})

func canParse(input string, expected *wat.Module) {
	result, err := wat.ParseString(input)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result).To(Equal(expected))
}
