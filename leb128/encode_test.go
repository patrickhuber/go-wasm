package leb128_test

import (
	"bufio"
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/leb128"
)

var _ = Describe("Encode", func() {
	DescribeTable("LebUint128",
		func(b []byte, value uint32) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			_, err := leb128.Encode(writer, value)
			Expect(err).To(BeNil())
			Expect(writer.Flush()).To(BeNil())
			Expect(buf.Bytes()).To(Equal(b))
		},
		Entry("one byte", []byte{0x08}, uint32(8)),
		Entry("two bytes", []byte{0x80, 0x7f}, uint32(16256)),
		Entry("three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)),
		Entry("five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)))
})
