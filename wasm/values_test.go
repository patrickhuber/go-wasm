package wasm_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/patrickhuber/go-wasm/wasm"
	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	buf := []byte("\x01C")
	count := 2
	expected := "C"

	t.Run("slice", func(t *testing.T) {
		value, read, err := wasm.DecodeUtf8String(buf)
		require.Nil(t, err)
		require.Equal(t, read, count)
		require.Equal(t, value, expected)
	})

	t.Run("reader", func(t *testing.T) {
		r := bytes.NewReader(buf)
		reader := bufio.NewReader(r)
		value, read, err := wasm.ReadUtf8String(reader)
		require.Nil(t, err)
		require.Equal(t, read, count)
		require.Equal(t, value, expected)
	})
}
