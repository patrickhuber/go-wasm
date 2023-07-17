package wit_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-wasm/wit"
	"github.com/patrickhuber/go-wasm/wit/token"
	"github.com/stretchr/testify/require"
)

func TestLexerUnit(t *testing.T) {
	type test struct {
		name      string
		str       string
		tokenType token.TokenType
	}
	tests := []test{
		{"line_comment", "// this is a comment line", token.LineComment},
		{"block_comment", "/* this is a comment block */", token.BlockComment},
		{"whitespace", "\f\t ", token.Whitespace},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lex := wit.NewLexer(strings.NewReader(test.str))
			tok, err := lex.Next()
			require.Nil(t, err)
			require.NotNil(t, tok)
			require.Equal(t, test.str, tok.Capture)
			require.Equal(t, test.tokenType, tok.Type)
		})
	}
}
