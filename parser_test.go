package wasm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm"
)

var _ = Describe("Parser", func() {
	Context("can parse", func() {
		It("module", func() {
			canParse("(module)", &wasm.Module{})
		})
		It("function", func() {
			canParse(`
			(module
				(func (param $lhs i32) (param $rhs i32) (result i32)
				  local.get $lhs
				  local.get $rhs
				  i32.add))`,
				&wasm.Module{
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
				})
		})
	})
})

func canParse(input string, expected *wasm.Module) {
	result, err := wasm.ParseString(input)
	Expect(err).To(BeNil())
	Expect(result).ToNot(BeNil())
	Expect(result).To(Equal(expected))
}
