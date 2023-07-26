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
			lexer := lex.New([]rune(test.str))
			tok, err := lexer.Next()
			require.NoError(t, err)
			require.NotNil(t, tok)
			require.Equal(t, test.str, string(tok.Runes))
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
			{Runes: []rune("variant"), Type: token.Variant},
			{Runes: []rune(" "), Type: token.Whitespace},
			{Runes: []rune("option"), Type: token.Option},
			{Runes: []rune(" "), Type: token.Whitespace},
			{Runes: []rune("{"), Type: token.OpenBrace},
			{Runes: []rune(" "), Type: token.Whitespace},
			{Runes: []rune("none"), Type: token.Id},
			{Runes: []rune(","), Type: token.Comma},
			{Runes: []rune(" "), Type: token.Whitespace},
			{Runes: []rune("some"), Type: token.Id},
			{Runes: []rune("("), Type: token.OpenParen},
			{Runes: []rune("ty"), Type: token.Id},
			{Runes: []rune(")"), Type: token.CloseParen},
			{Runes: []rune(","), Type: token.Comma},
			{Runes: []rune(" "), Type: token.Whitespace},
			{Runes: []rune("}"), Type: token.CloseBrace},
		}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			lex := lex.New([]rune(test.str))
			for i, tok := range test.tokens {
				actual, err := lex.Next()
				require.NoError(t, err, "%d", i)
				require.Equal(t, tok.Type, actual.Type, "%d", i)
				require.Equal(t, tok.Runes, actual.Runes, "%d", i)
			}
		})
	}
	t.Run("peek", func(t *testing.T) {
		str := "variant option { none, some(ty), }"
		lex := lex.New([]rune(str))
		tok, err := lex.Peek()
		for i := 0; i < 10; i++ {
			require.NoError(t, err)
			require.Equal(t, token.Variant, tok.Type)
		}
		tok, err = lex.Next()
		require.NoError(t, err)
		require.Equal(t, tok.Type, token.Variant)
		require.Equal(t, "variant", string(tok.Runes))

		tok, err = lex.Peek()
		require.NoError(t, err)
		require.Equal(t, tok.Type, token.Whitespace)
		require.Equal(t, " ", string(tok.Runes))
	})
}
