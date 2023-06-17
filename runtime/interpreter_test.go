package runtime_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/go-wasm/runtime"
	"github.com/stretchr/testify/require"
)

func TestI32Add(t *testing.T) {
	body := []byte{}
	i := runtime.NewInterpreter()
	err := i.Run(bytes.NewReader(body))
	require.Nil(t, err)
}
