package wasm_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-wasm/wasm"
)

var _ = Describe("Header", func() {

	Describe("Read", func() {
		It("can read module header", func() {
			buf := []byte("\x00asm\x01\x00\x00\x00")
			reader := wasm.NewHeaderReader(bytes.NewBuffer(buf))
			header, err := reader.Read()
			Expect(err).To(BeNil())
			Expect(header).ToNot(BeNil())
			Expect(*header).To(Equal(wasm.Header{
				Magic:   wasm.Magic,
				Version: wasm.Version1,
				Layer:   wasm.LayerCore,
			}))
		})
		It("can read component header", func() {
			buf := []byte("\x00asm\x0a\x00\x01\x00")
			reader := wasm.NewHeaderReader(bytes.NewBuffer(buf))
			header, err := reader.Read()
			Expect(err).To(BeNil())
			Expect(header).ToNot(BeNil())
			Expect(*header).To(Equal(wasm.Header{
				Magic:   wasm.Magic,
				Version: wasm.VersionExperimental,
				Layer:   wasm.LayerComponent,
			}))
		})
	})
	Describe("Write", func() {

	})
})
