package wit_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-wasm/wit"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		text := `package local:demo

interface console {
  log: func(arg: string)
}
`
		reader := strings.NewReader(text)
		node, err := wit.Parse(reader)
		require.Nil(t, err)
		require.NotNil(t, node)
	})
}
