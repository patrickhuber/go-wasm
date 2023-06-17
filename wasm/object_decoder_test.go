package wasm_test

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/patrickhuber/go-wasm/internal/to"
	"github.com/patrickhuber/go-wasm/wasm"
	"github.com/stretchr/testify/require"
)

func TestObjectDecoder(t *testing.T) {
	type test struct {
		name     string
		object   *wasm.Object
		filePath string
	}
	tests := []test{
		{"empty", &wasm.Object{
			Header: wasm.NewModuleHeader(),
			Module: &wasm.Module{},
		}, "../fixtures/empty/empty.wasm"},
		{"func", &wasm.Object{
			Header: wasm.NewModuleHeader(),
			Module: &wasm.Module{
				Sections: []wasm.Section{
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
					},
					{
						ID:   wasm.FuncSectionType,
						Size: 2,
						Function: &wasm.FunctionSection{
							Types: []uint32{0},
						},
					},
					{
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
					},
				},
			},
		}, "../fixtures/func/func.wasm"},
		{"component_empty", &wasm.Object{
			Header:    wasm.NewComponentHeader(),
			Component: &wasm.Component{},
		}, "../fixtures/component/empty.wasm"},
		{
			"component_name",
			&wasm.Object{
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
			},
			"../fixtures/component/name.wasm",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, filename, _, _ := runtime.Caller(0)
			filePath := path.Join(path.Dir(filename), test.filePath)

			f, err := os.Open(filePath)
			require.Nil(t, err)

			defer f.Close()
			reader := wasm.NewObjectDecoder(f)
			read, err := reader.Decode()
			require.Nil(t, err)
			require.Equal(t, read, test.object)
		})
	}
}
