package binary_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/go-wasm/binary"
	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/instruction"
	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		document *binary.Document
	}{
		{
			name: "empty",
			path: "../fixtures/empty/empty.wasm",
			document: &binary.Document{
				Preamble: &binary.Preamble{
					Magic:   binary.Magic,
					Version: binary.ModuleVersion,
					Layer:   0,
				},
				Directive: &binary.Module{},
			},
		},
		{
			name: "func",
			path: "../fixtures/func/func.wasm",
			document: &binary.Document{
				Preamble: &binary.Preamble{
					Magic:   binary.Magic,
					Version: binary.ModuleVersion,
					Layer:   0,
				},
				Directive: &binary.Module{
					Sections: []binary.Section{
						&binary.TypeSection{
							ID:   binary.TypeSectionID,
							Size: 4,
							Types: []*binary.FunctionType{
								{
									Parameters: []binary.ValueType{},
									Results:    []binary.ValueType{},
								},
							},
						},
						&binary.FunctionSection{
							ID:    binary.FunctionSectionID,
							Size:  2,
							Types: []uint32{0},
						},
						&binary.CodeSection{
							ID:   binary.CodeSectionID,
							Size: 4,
							Codes: []*binary.Code{
								{
									Size:   2,
									Locals: []binary.Local{},
									Expression: []instruction.Instruction{
										instruction.End{},
									},
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
			document: &binary.Document{
				Preamble: &binary.Preamble{
					Magic:   binary.Magic,
					Version: binary.ModuleVersion,
					Layer:   0,
				},
				Directive: &binary.Module{
					Sections: []binary.Section{
						&binary.TypeSection{
							ID:   binary.TypeSectionID,
							Size: 7,
							Types: []*binary.FunctionType{
								{
									Parameters: []binary.ValueType{
										binary.I32,
										binary.I32,
									},
									Results: []binary.ValueType{
										binary.I32,
									},
								},
							},
						},
						&binary.FunctionSection{
							ID:    binary.FunctionSectionID,
							Size:  2,
							Types: []uint32{0},
						},
						&binary.CodeSection{
							ID:   binary.CodeSectionID,
							Size: 9,
							Codes: []*binary.Code{
								{
									Size:   7,
									Locals: []binary.Local{},
									Expression: []instruction.Instruction{
										instruction.LocalGet{
											Index: indicies.Local(0),
										},
										instruction.LocalGet{
											Index: indicies.Local(1),
										},
										instruction.I32Add{},
										instruction.End{},
									},
								},
							},
						},
					},
				},
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
