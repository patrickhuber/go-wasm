package wasm_test

import (
	"bufio"
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wasm"
)

var _ = Describe("Values", func() {
	It("can read string from slice", func() {

		value, read, err := wasm.DecodeUtf8String([]byte("\x01C"))
		Expect(err).To(BeNil())
		Expect(read).To(Equal(2))
		Expect(value).To(Equal("C"))
	})
	It("can read string from reader", func() {
		r := bytes.NewBufferString("\x01C")
		reader := bufio.NewReader(r)
		value, read, err := wasm.ReadUtf8String(reader)
		Expect(err).To(BeNil())
		Expect(read).To(Equal(2))
		Expect(value).To(Equal("C"))
	})
})
