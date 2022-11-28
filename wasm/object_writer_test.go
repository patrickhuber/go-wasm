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
		object := &wasm.Object{
			Header: wasm.NewModuleHeader(),
			Module: &wasm.Module{},
		}
		var buffer bytes.Buffer
		writer := wasm.NewObjectWriter(&buffer)
		err := writer.Write(object)
		Expect(err).To(BeNil())
		Expect(buffer.Bytes()).To(Equal(expected))
	})
})
