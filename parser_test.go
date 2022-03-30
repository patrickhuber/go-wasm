package wasm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm"
	"github.com/patrickhuber/go-wasm/builder"
)

var _ = Describe("Parser", func() {
	Describe("Parse", func() {
		It("module", func() {
			canParse("(module)", &wasm.Module{})
		})
		It("memory", func() {
			b := builder.NewModule(func(s builder.Section) {
				s.Memory(func(m builder.Memory) {
					m.Limits(1)
				})
				s.Function(func(f builder.Function) {})
			})
			canParse("(module (memory 1) (func))", b.Build())
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
								p.Parameter(wasm.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 parameter", func() {
				canParse("(module (func (param i64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(wasm.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 parameter", func() {
				canParse("(module (func (param f32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(wasm.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 parameter", func() {
				canParse("(module (func (param f64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Parameters(func(p builder.Parameters) {
								p.Parameter(wasm.F64)
							})
						})
					}).Build())
			})
			It("can parse i32 result", func() {
				canParse("(module (func (result i32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(wasm.I32)
							})
						})
					}).Build())
			})
			It("can parse i64 result", func() {
				canParse("(module (func (result i64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(wasm.I64)
							})
						})
					}).Build())
			})
			It("can parse f32 result", func() {
				canParse("(module (func (result f32)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(wasm.F32)
							})
						})
					}).Build())
			})
			It("can parse f64 result", func() {
				canParse("(module (func (result f64)))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(wasm.F64)
							})
						})
					}).Build())
			})
			It("can parse mutiple results", func() {
				canParse("(module (func (result i64) (result i64) ))", builder.NewModule(
					func(s builder.Section) {
						s.Function(func(f builder.Function) {
							f.Results(func(p builder.Results) {
								p.Result(wasm.I64)
								p.Result(wasm.I64)
							})
						})
					}).Build())
			})
		})
		It("function", func() {
			builder := builder.NewModule(func(s builder.Section) {
				s.Function(func(f builder.Function) {
					f.Parameters(func(p builder.Parameters) {
						p.Parameter(wasm.I32).ID("$lhs")
						p.Parameter(wasm.I32).ID("$rhs")
					})
					f.Results(func(r builder.Results) {
						r.Result(wasm.I32)
					})
					f.Instructions(func(i builder.Instructions) {
						i.Local(wasm.LocalGet).ID("$lhs")
						i.Local(wasm.LocalGet).ID("$rhs")
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

func canParse(input string, expected *wasm.Module) {
	result, err := wasm.ParseString(input)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result).To(Equal(expected))
}
