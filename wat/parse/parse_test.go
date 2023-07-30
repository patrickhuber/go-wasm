package parse_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/go-wasm/wat/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"
	"github.com/patrickhuber/go-wasm/wat/parse"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("empty_module", func(t *testing.T) {
		lexer := lex.New("(module)")
		n, err := parse.Parse(lexer)
		require.NoError(t, err)
		require.NotNil(t, n)
		_, ok := n.(*ast.Module)
		require.True(t, ok)
	})
	t.Run("empty_component", func(t *testing.T) {
		lexer := lex.New("(component)")
		n, err := parse.Parse(lexer)
		require.NoError(t, err)
		require.NotNil(t, n)
		_, ok := n.(*ast.Component)
		require.True(t, ok)
	})
	t.Run("add", func(t *testing.T) {
		file := "../../fixtures/add/add.wat"
		content, err := os.ReadFile(file)
		require.NoError(t, err)
		lexer := lex.New(string(content))
		n, err := parse.Parse(lexer)
		require.NoError(t, err)
		require.NotNil(t, n)
	})
}
