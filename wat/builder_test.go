package wat_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-wasm/to"
	"github.com/patrickhuber/go-wasm/wat"
)

var _ = Describe("Builder", func() {
	It("can create empty module", func() {
		builder := wat.NewModule(func(s wat.SectionBuilder) {
			s.Function(func(f wat.FunctionBuilder) {
				f.Parameters(func(p wat.ParametersBuilder) {
					p.Parameter(wat.I32).ID("$lhs")
				})
			})
		})
		module := builder.Build()
		Expect(module).ToNot(BeNil())
	})
	It("can build memory", func() {
		builder := wat.NewModule(func(s wat.SectionBuilder) {
			s.Memory(func(m wat.MemoryBuilder) {
				m.Limits(1)
			})
		})
		module := &wat.Module{
			Memory: []wat.Section{
				{
					Memory: &wat.Memory{
						Limits: wat.Limits{
							Min: 1,
						},
					},
				},
			},
		}
		Expect(builder.Build()).To(Equal(module))
	})
	It("can build function", func() {
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
		module := &wat.Module{
			Functions: []wat.Section{
				{
					Function: &wat.Function{
						ID: nil,
						Parameters: []wat.Parameter{
							{
								ID:   to.Pointer(wat.Identifier("$lhs")),
								Type: wat.I32,
							},
							{
								ID:   to.Pointer(wat.Identifier("$rhs")),
								Type: wat.I32,
							},
						},
						Results: []wat.Result{
							{
								Type: wat.I32,
							},
						},
						Instructions: []wat.Instruction{
							{
								Plain: &wat.Plain{
									Local: &wat.LocalInstruction{
										Operation: wat.LocalGet,
										ID:        to.Pointer(wat.Identifier("$lhs")),
									},
								},
							},
							{
								Plain: &wat.Plain{
									Local: &wat.LocalInstruction{
										Operation: wat.LocalGet,
										ID:        to.Pointer(wat.Identifier("$rhs")),
									},
								},
							},
							{
								Plain: &wat.Plain{
									I32: &wat.I32Instruction{
										Operation: wat.BinaryOperationAdd,
									},
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
