package binary_test

import (
	"bytes"
	"os"
	"slices"
	"testing"

	"github.com/patrickhuber/go-wasm/binary"
	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/instruction"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		document *binary.Document
	}{
		{
			name: "add",
			path: "../fixtures/add/add.wasm",
			document: &binary.Document{
				Preamble: &binary.Preamble{
					Magic:   [4]byte{0x00, 0x61, 0x73, 0x6d},
					Version: 1,
					Layer:   0,
				},
				Root: &binary.Module{
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

func TestReadPreamble(t *testing.T) {
	type test struct {
		name     string
		preamble *binary.Preamble
		buffer   []byte
	}
	tests := []test{
		{
			name: "component",
			preamble: &binary.Preamble{
				Magic:   [4]byte{0x00, 0x61, 0x73, 0x6d},
				Version: 13,
				Layer:   1,
			},
			buffer: []byte{
				// magic
				0x00, 0x61, 0x73, 0x6d,
				// version
				0x0d, 0x00,
				// layer
				0x01, 0x00},
		},
		{
			name: "module",
			preamble: &binary.Preamble{
				Magic:   [4]byte{0x00, 0x61, 0x73, 0x6d},
				Version: 1,
				Layer:   0,
			},
			buffer: []byte{
				// magic
				0x00, 0x61, 0x73, 0x6d,
				// version
				0x01, 0x00,
				// layer
				0x00, 0x00},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reader := bytes.NewBuffer(test.buffer)
			preamble, err := binary.ReadPreamble(reader)
			if err != nil {
				t.Fatal(err)
			}
			if preamble == nil {
				t.Fatalf("preamble is nil")
			}
			if !slices.Equal(preamble.Magic[0:], test.preamble.Magic[0:]) {
				t.Fatalf("invalid magic number %v", preamble.Magic)
			}
			if preamble.Version != test.preamble.Version {
				t.Fatalf("expected version 13, found %v", preamble.Version)
			}
			if preamble.Layer != test.preamble.Layer {
				t.Fatalf("expected layer 1, found %v", preamble.Layer)
			}
		})

	}
}
