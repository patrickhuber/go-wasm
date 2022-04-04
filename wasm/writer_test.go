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
		var buffer bytes.Buffer
		writer := wasm.NewWriter(&buffer)
		err := writer.Write(module)
		Expect(err).To(BeNil())
		Expect(buffer.Bytes()).To(Equal(expected))
	})
})
