package wasm_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wasm"
)

var _ = Describe("Writer", func() {
	It("can write empty module", func() {
		expected := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
		module := &wasm.Module{
			Magic:   wasm.Magic,
			Version: wasm.Version,
		}
		buffer := &bytes.Buffer{}
		writer := wasm.NewWriter(buffer)
		err := writer.Write(module)
		Expect(err).To(BeNil())
		Expect(buffer.Bytes()).To(Equal(expected))
	})
	DescribeTable("LebU128",
		func(b []byte, value uint32) {
			buffer := &bytes.Buffer{}
			writer := wasm.NewWriter(buffer)
			err := writer.WriteLebU128(value)
			Expect(err).To(BeNil())
			Expect(buffer.Bytes()).To(Equal(b))
		},
		Entry("one byte", []byte{0x08}, uint32(8)),
		Entry("two bytes", []byte{0x80, 0x7f}, uint32(16256)),
		Entry("three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)),
		Entry("five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)))
})
