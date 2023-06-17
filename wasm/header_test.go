package wasm_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/go-wasm/wasm"
	"github.com/stretchr/testify/require"
)

func TestCanReadModuleHeader(t *testing.T) {
	buf := []byte("\x00asm\x01\x00\x00\x00")
	reader := wasm.NewHeaderReader(bytes.NewBuffer(buf))
	header, err := reader.Read()
	require.Nil(t, err)
	require.NotNil(t, header)
	require.Equal(t, *header, wasm.Header{
		Magic:   wasm.Magic,
		Version: wasm.Version1,
		Layer:   wasm.LayerCore,
	})
}

func TestCanReadComponentHeader(t *testing.T) {
	buf := []byte("\x00asm\x0a\x00\x01\x00")
	reader := wasm.NewHeaderReader(bytes.NewBuffer(buf))
	header, err := reader.Read()
	require.Nil(t, err)
	require.NotNil(t, header)
	require.Equal(t, *header, wasm.Header{
		Magic:   wasm.Magic,
		Version: wasm.VersionExperimental,
		Layer:   wasm.LayerComponent,
	})
}
