package lex_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/wit/lex"
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
		{"variant", "variant", token.Variant},
		{"open_brace", "{", token.OpenBrace},
		{"close_brace", "}", token.CloseBrace},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := lex.New(test.str)
			tok, err := lexer.Next()
			require.NoError(t, err)
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
			{Capture: "variant", Type: token.Variant},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "option", Type: token.Option},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "{", Type: token.OpenBrace},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "none", Type: token.Id},
			{Capture: ",", Type: token.Comma},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "some", Type: token.Id},
			{Capture: "(", Type: token.OpenParen},
			{Capture: "ty", Type: token.Id},
			{Capture: ")", Type: token.CloseParen},
			{Capture: ",", Type: token.Comma},
			{Capture: " ", Type: token.Whitespace},
			{Capture: "}", Type: token.CloseBrace},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			lex := lex.New(test.str)
			for i, tok := range test.tokens {
				actual, err := lex.Next()
				require.NoError(t, err, "%d", i)
				require.Equal(t, tok.Type, actual.Type, "%d", i)
				require.Equal(t, tok.Capture, actual.Capture, "%d", i)
			}
		})
	}
	t.Run("peek", func(t *testing.T) {
		str := "variant option { none, some(ty), }"
		lex := lex.New(str)
		tok, err := lex.Peek()
		for i := 0; i < 10; i++ {
			require.NoError(t, err)
			require.Equal(t, token.Variant, tok.Type)
		}
		tok, err = lex.Next()
		require.NoError(t, err)
		require.Equal(t, tok.Type, token.Variant)
		require.Equal(t, "variant", string(tok.Capture))

		tok, err = lex.Peek()
		require.NoError(t, err)
		require.Equal(t, tok.Type, token.Whitespace)
		require.Equal(t, " ", string(tok.Capture))
	})
}
