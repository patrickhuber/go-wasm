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
	type test struct {
		name     string
		text     string
		expected ast.Ast
	}
	tests := []test{
		{"empty_module", "(module)", &ast.Module{}},
		{"empty_component", "(component)", &ast.Component{}},
		{"export_function", `(module (func (export "add")))`, &ast.Module{Functions: []ast.Function{{Exports: []ast.InlineExport{{Name: `"add"`}}}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := lex.New(test.text)
			n, err := parse.Parse(lexer)
			require.NoError(t, err)
			require.NotNil(t, n)
			require.EqualValues(t, test.expected, n)
		})
	}
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
