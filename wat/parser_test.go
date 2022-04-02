package wat_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/builder"
	"github.com/patrickhuber/go-wasm/model"
	"github.com/patrickhuber/go-wasm/wat"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		It("module", func() {
			canParse("(module)", &model.Module{})
		})
		It("memory", func() {
			canParse("(module (memory 1) (func))", builder.NewModule(func(s builder.Section) {
				s.Memory(func(m builder.Memory) {
					m.Limits(1)
				})
				s.Function(func(f builder.Function) {})
			}).Build())
		})
		Describe("function", func() {
			It("can parse function alias", func() {
				canParse("(module (func $alias ))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.ID("$alias")
						})
					}).Build())
			})
			It("can parse i32 parameter", func() {
				canParse("(module (func (param i32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(model.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 parameter", func() {
				canParse("(module (func (param i64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(model.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 parameter", func() {
				canParse("(module (func (param f32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(model.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 parameter", func() {
				canParse("(module (func (param f64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(model.F64)
							})
						})
					}).Build())
			})
			It("can parse i32 result", func() {
				canParse("(module (func (result i32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(model.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 result", func() {
				canParse("(module (func (result i64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(model.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 result", func() {
				canParse("(module (func (result f32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(model.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 result", func() {
				canParse("(module (func (result f64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(model.F64)
							})
						})
					}).Build())
			})
			It("can parse mutiple results", func() {
				canParse("(module (func (result i64) (result i64) ))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(model.I64)
								p.Result(model.I64)
							})
						})
					}).Build())
			})
		})
		It("function", func() {
			builder := builder.NewModule(func(s builder.Section) {
				s.Function(func(f builder.Function) {
					f.Parameters(func(p builder.Parameters) {
						p.Parameter(model.I32).ID("$lhs")
						p.Parameter(model.I32).ID("$rhs")
					})
					f.Results(func(r builder.Results) {
						r.Result(model.I32)
					})
					f.Instructions(func(i builder.Instructions) {
						i.Local(model.LocalGet).ID("$lhs")
						i.Local(model.LocalGet).ID("$rhs")
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

func canParse(input string, expected *model.Module) {
	result, err := wat.ParseString(input)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result).To(Equal(expected))
}
