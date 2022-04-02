package wasm_test

import (
	"bytes"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wasm"
)

var cases = []struct {
	v uint32
	b []byte
}{
	{b: []byte{0x08}, v: 8},
	{b: []byte{0x80, 0x7f}, v: 16256},
	{b: []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, v: 2141192192},
}

var _ = Describe("Reader", func() {
	DescribeTable("LebUint128",
		func(b []byte, value uint32) {
			reader := wasm.NewReader(bytes.NewBuffer(b))
			result, err := reader.ReadLebU128()
			Expect(err).To(BeNil())
			Expect(result).To(Equal(value))

		},
		Entry("one byte", []byte{0x08}, uint32(8)),
		Entry("two bytes", []byte{0x80, 0x7f}, uint32(16256)),
		Entry("three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)),
		Entry("five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)))
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
