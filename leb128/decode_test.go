package leb128_test

import (
	"bufio"
	"bytes"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/leb128"
	"github.com/stretchr/testify/require"
)

var _ = Describe("Decode", func() {
	DescribeTable("LebUint128Slice",
		func(b []byte, value uint32) {
			result, _, err := leb128.DecodeSlice(b)
			Expect(err).To(BeNil())
			Expect(value).To(Equal(result))
		},
		Entry("one byte", []byte{0x08}, uint32(8)),
		Entry("two bytes", []byte{0x80, 0x7f}, uint32(16256)),
		Entry("three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)),
		Entry("five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)))

})

func TestLebUint128DecodeSlice(t *testing.T) {
	type test struct {
		name  string
		buf   []byte
		value uint32
	}
	tests := []test{
		{"one byte", []byte{0x08}, uint32(8)},
		{"two bytes", []byte{0x80, 0x7f}, uint32(16256)},
		{"three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)},
		{"five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, _, err := leb128.DecodeSlice(test.buf)
			require.Nil(t, err)
			require.Equal(t, test.value, result)
		})
	}
}

func TestLebUint128Decode(t *testing.T) {
	type test struct {
		name  string
		buf   []byte
		value uint32
	}
	tests := []test{
		{"one byte", []byte{0x08}, uint32(8)},
		{"two bytes", []byte{0x80, 0x7f}, uint32(16256)},
		{"three bytes", []byte{0xE5, 0x8E, 0x26}, uint32(624485)},
		{"five bytes", []byte{0x80, 0x80, 0x80, 0xfd, 0x07}, uint32(2141192192)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := bufio.NewReader(bytes.NewReader(test.buf))
			result, _, err := leb128.Decode(r)
			require.Nil(t, err)
			require.Equal(t, test.value, result)
		})
	}
}
