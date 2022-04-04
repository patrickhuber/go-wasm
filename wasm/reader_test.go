package wasm_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wasm"
)

var _ = Describe("Reader", func() {
	It("can read empty module", func() {
		m := &wasm.Module{
			Magic:   wasm.Magic,
			Version: wasm.Version,
		}
		equal("fixtures/empty.wasm", m)
	})
	It("can read empty func", func() {
		m := &wasm.Module{
			Magic:   wasm.Magic,
			Version: wasm.Version,
			Types: []wasm.Section{
				{
					ID:   wasm.TypeSectionType,
					Size: 4,
					Type: &wasm.TypeSection{
						Types: []wasm.Type{
							{
								Parameters: &wasm.ResultType{
									Values: []*wasm.ValueType{},
								},
								Results: &wasm.ResultType{
									Values: []*wasm.ValueType{},
								},
							},
						},
					},
				}},
			Functions: []wasm.Section{
				{
					ID:   wasm.FuncSectionType,
					Size: 2,
					Function: &wasm.FunctionSection{
						Types: []uint32{0},
					},
				}},
			Codes: []wasm.Section{{
				ID:   wasm.CodeSectionType,
				Size: 4,
				Code: &wasm.CodeSection{
					Codes: []wasm.Code{
						{
							Size: 2,
							Expression: []wasm.Instruction{
								{
									OpCode: wasm.End,
								},
							},
						},
					},
				},
			}},
		}
		equal("fixtures/func.wasm", m)
	})
})

func equal(file string, module *wasm.Module) {
	f, err := os.Open(file)
	Expect(err).To(BeNil())
	defer f.Close()
	reader := wasm.NewReader(f)
	read, err := reader.Read()
	Expect(err).To(BeNil())
	Expect(read).To(Equal(module))
}
