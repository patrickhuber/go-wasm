package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-wasm"
	"github.com/patrickhuber/go-wasm/builder"
)

var _ = Describe("Builder", func() {
	It("can create empty module", func() {
		builder := builder.NewModule(func(s builder.Section) {
			s.Function(func(f builder.Function) {
				f.Parameters(func(p builder.Parameters) {
					p.Parameter(wasm.I32).ID("$lhs")
				})
			})
		})
		module := builder.Build()
		Expect(module).ToNot(BeNil())
	})
	It("can build memory", func() {
		builder := builder.NewModule(func(s builder.Section) {
			s.Memory(func(m builder.Memory) {
				m.Limits(1)
			})
		})
		module := &wasm.Module{
			Memory: []wasm.Memory{
				{
					Limits: wasm.Limits{
						Min: 1,
					},
				},
			},
		}
		Expect(builder.Build()).To(Equal(module))
	})
	It("can build function", func() {
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
		module := &wasm.Module{
			Functions: []wasm.Function{
				{
					ID: nil,
					Parameters: []wasm.Parameter{
						{
							ID:   wasm.Pointer(wasm.Identifier("$lhs")),
							Type: wasm.I32,
						},
						{
							ID:   wasm.Pointer(wasm.Identifier("$rhs")),
							Type: wasm.I32,
						},
					},
					Results: []wasm.Result{
						{
							Type: wasm.I32,
						},
					},
					Instructions: []wasm.Instruction{
						{
							Plain: &wasm.Plain{
								Local: &wasm.LocalInstruction{
									Operation: wasm.LocalGet,
									ID:        wasm.Pointer(wasm.Identifier("$lhs")),
								},
							},
						},
						{
							Plain: &wasm.Plain{
								Local: &wasm.LocalInstruction{
									Operation: wasm.LocalGet,
									ID:        wasm.Pointer(wasm.Identifier("$rhs")),
								},
							},
						},
						{
							Plain: &wasm.Plain{
								I32: &wasm.I32Instruction{
									Operation: wasm.BinaryOperationAdd,
								},
							},
						},
					},
				},
			},
		}
		Expect(builder.Build()).To(Equal(module))
	})
})
