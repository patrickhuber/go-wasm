package wasm_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/go-wasm/wasm"
	"github.com/stretchr/testify/require"
)

func TestWrite(t *testing.T) {
	type test struct {
		name     string
		object   *wasm.Object
		expected []byte
	}
	tests := []test{
		{
			"empty",
			&wasm.Object{
				Header: wasm.NewModuleHeader(),
				Module: &wasm.Module{},
			},
			[]byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			writer := wasm.NewObjectEncoder(&buffer)
			err := writer.Encode(test.object)
			require.Nil(t, err)
			require.Equal(t, buffer.Bytes(), test.expected)
		})
	}
}
