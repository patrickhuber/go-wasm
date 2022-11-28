package wasm_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/to"
	"github.com/patrickhuber/go-wasm/wasm"
)

var _ = Describe("ObjectReader", func() {
	Describe("Module", func() {
		It("can read empty module", func() {
			o := &wasm.Object{
				Header: wasm.NewModuleHeader(),
				Module: &wasm.Module{},
			}
			equal("fixtures/empty.wasm", o)
		})
		It("can read empty module func", func() {
			o := &wasm.Object{
				Header: wasm.NewModuleHeader(),
				Module: &wasm.Module{
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
				},
			}
			equal("fixtures/func.wasm", o)
		})
	})
	Describe("Component", func() {

		It("can read empty component", func() {
			object := &wasm.Object{
				Header:    wasm.NewComponentHeader(),
				Component: &wasm.Component{},
			}
			equal("fixtures/empty_component.wasm", object)
		})
		It("can read named component", func() {
			object := &wasm.Object{
				Header: wasm.NewComponentHeader(),
				Component: &wasm.Component{
					Custom: []wasm.Section{
						{
							ID:   wasm.CustomSectionType,
							Size: 19,
							Custom: &wasm.CustomSection{
								Name: &wasm.NameSection{
									Key:  "component-name",
									Name: to.Pointer("C"),
								},
							},
						},
					},
				},
			}
			equal("fixtures/name_component.wasm", object)
		})
	})

})

func equal(file string, object *wasm.Object) {
	f, err := os.Open(file)
	Expect(err).To(BeNil())
	defer f.Close()
	reader := wasm.NewObjectReader(f)
	read, err := reader.Read()
	Expect(err).To(BeNil())
	Expect(read).To(Equal(object))
}
