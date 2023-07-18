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
		{"whitespace", "\f\t \r\n", token.Whitespace},
		{"string", "variant", token.String},
		{"open_brace", "{", token.OpenBrace},
		{"close_brace", "}", token.CloseBrace},
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

func TestLexer(t *testing.T) {
	type test struct {
		name   string
		str    string
		tokens []token.Token
	}
	tests := []test{
		{"variant", "variant option { none, some(ty), }", []token.Token{
			{Capture: "variant", Type: token.String},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "option", Type: token.String},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "{", Type: token.OpenBrace},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "none", Type: token.String},
			{Capture: ",", Type: token.Comma},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "some", Type: token.String},
			{Capture: "(", Type: token.OpenParen},
			{Capture: "ty", Type: token.String},
			{Capture: ")", Type: token.CloseParen},
			{Capture: ",", Type: token.Comma},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "}", Type: token.CloseBrace},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			lex := wit.NewLexer(strings.NewReader(test.str))
			for i, tok := range test.tokens {
				actual, err := lex.Next()
				require.Nil(t, err)
				require.Equal(t, i, i) // needed for conditional breakpoint
				require.Equal(t, tok.Type, actual.Type)
				require.Equal(t, tok.Capture, actual.Capture)
			}
		})
	}
}
