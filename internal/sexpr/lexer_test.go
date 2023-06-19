package sexpr_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/internal/sexpr"
	"github.com/stretchr/testify/require"
)

func TestCanTokenizeModule(t *testing.T) {
	CanTokenize(t, "(module)",
		sexpr.OpenParen,
		sexpr.String,
		sexpr.CloseParen)
}

func TestCanTokenizeWhitespace(t *testing.T) {
	CanTokenize(t, " (module ) ",
		sexpr.Whitespace,
		sexpr.OpenParen,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.CloseParen,
		sexpr.Whitespace)
}

func TestCanTokenizeFunction(t *testing.T) {
	CanTokenize(t, `
	( module 
		(func $add (param $lhs i32) (param $rhs i32) (result i32)
			local.get $lhs
			local.get $rhs
			i32.add))`,
		sexpr.Whitespace,
		sexpr.OpenParen,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,

		sexpr.OpenParen,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,

		sexpr.OpenParen,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.CloseParen,
		sexpr.Whitespace,

		sexpr.OpenParen,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.CloseParen,
		sexpr.Whitespace,

		sexpr.OpenParen,
		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.CloseParen,
		sexpr.Whitespace,

		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,

		sexpr.String,
		sexpr.Whitespace,
		sexpr.String,
		sexpr.Whitespace,

		sexpr.String,
		sexpr.CloseParen,

		sexpr.CloseParen,
	)
}

func CanTokenize(t *testing.T, input string, sequence ...sexpr.TokenType) {
	lexer := sexpr.NewLexer(input)
	for i, item := range sequence {
		token, err := lexer.Next()
		require.Nil(t, err)
		require.Equal(t, token.Type, item, "at %d", i)
	}
	token, err := lexer.Next()
	require.Nil(t, err)
	require.Equal(t, token.Type, sexpr.EndOfStream, "at %d", len(sequence)-1)
}
