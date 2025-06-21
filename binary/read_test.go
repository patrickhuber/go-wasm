package binary_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/binary"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		document *api.Document
	}{
		{
			name: "empty",
			path: "../fixtures/empty/empty.wasm",
			document: &api.Document{
				Preamble: api.Preamble{
					Version: api.ModuleVersion,
					Layer:   0,
				},
				Directive: &api.Module{},
			},
		},
		{
			name: "func",
			path: "../fixtures/func/func.wasm",
			document: &api.Document{
				Preamble: api.Preamble{
					Version: api.ModuleVersion,
					Layer:   0,
				},
				Directive: &api.Module{
					Types: []*api.FuncType{
						&api.FuncType{
							Parameters: api.ResultType{
								Types: []api.ValType{},
							},
							Returns: api.ResultType{
								Types: []api.ValType{},
							},
						},
					},
					Funcs: []*api.Func{
						&api.Func{
							Locals: []api.ValType{},
							Body: &api.Expression{
								Instructions: []api.Instruction{
									&api.End{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "add",
			path: "../fixtures/add/add.wasm",
			document: &api.Document{
				Preamble: api.Preamble{
					Version: api.ModuleVersion,
					Layer:   0,
				},
				Directive: &api.Module{
					Types: []*api.FuncType{
						&api.FuncType{
							Parameters: api.ResultType{
								Types: []api.ValType{
									&api.I32Type{},
									&api.I32Type{},
								},
							},
							Returns: api.ResultType{
								Types: []api.ValType{
									&api.I32Type{},
								},
							},
						},
					},
					Funcs: []*api.Func{
						&api.Func{
							Type:   api.TypeIndex(0),
							Locals: []api.ValType{},
							Body: &api.Expression{
								Instructions: []api.Instruction{
									&api.LocalGet{
										Index: api.LocalIndex(0),
									},
									&api.LocalGet{
										Index: api.LocalIndex(1),
									},
									&api.I32Add{},
									&api.End{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "component",
			path: "../fixtures/component/empty.wasm",
			document: &api.Document{
				Preamble: api.Preamble{
					Version: binary.ComponentVersion,
					Layer:   1,
				},
				Directive: &api.Component{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader, err := os.Open(test.path)
			if err != nil {
				t.Fatal(err)
			}
			d, err := binary.Read(reader)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, test.document, d)
		})
	}
}
