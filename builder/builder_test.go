package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-wasm/builder"
	"github.com/patrickhuber/go-wasm/model"
	"github.com/patrickhuber/go-wasm/to"
)

var _ = Describe("Builder", func() {
	It("can create empty module", func() {
		builder := builder.NewModule(func(s builder.Section) {
			s.Function(func(f builder.Function) {
				f.Parameters(func(p builder.Parameters) {
					p.Parameter(model.I32).ID("$lhs")
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
		module := &model.Module{
			Memory: []model.Memory{
				{
					Limits: model.Limits{
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
		module := &model.Module{
			Functions: []model.Function{
				{
					ID: nil,
					Parameters: []model.Parameter{
						{
							ID:   to.Pointer(model.Identifier("$lhs")),
							Type: model.I32,
						},
						{
							ID:   to.Pointer(model.Identifier("$rhs")),
							Type: model.I32,
						},
					},
					Results: []model.Result{
						{
							Type: model.I32,
						},
					},
					Instructions: []model.Instruction{
						{
							Plain: &model.Plain{
								Local: &model.LocalInstruction{
									Operation: model.LocalGet,
									ID:        to.Pointer(model.Identifier("$lhs")),
								},
							},
						},
						{
							Plain: &model.Plain{
								Local: &model.LocalInstruction{
									Operation: model.LocalGet,
									ID:        to.Pointer(model.Identifier("$rhs")),
								},
							},
						},
						{
							Plain: &model.Plain{
								I32: &model.I32Instruction{
									Operation: model.BinaryOperationAdd,
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
