package lex_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/wat/lex"
	"github.com/patrickhuber/go-wasm/wat/token"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	t.Run("module", func(t *testing.T) {
		CanTokenize(t, "(module)",
			token.OpenParen,
			token.Reserved,
			token.CloseParen)
	})
	t.Run("module_whitespace", func(t *testing.T) {
		CanTokenize(t, " (module ) ",
			token.Whitespace,
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.CloseParen,
			token.Whitespace)
	})

	t.Run("comments", func(t *testing.T) {
		CanTokenize(t, `
		;; Line Comment
		( module 
			(memory (; in the middle ;) 1)
		)`,
			token.Whitespace,
			token.LineComment,
			token.Whitespace,
			token.OpenParen,
			token.Whitespace,
			token.Reserved,
			token.Whitespace,
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.BlockComment,
			token.Whitespace,
			token.Integer,
			token.CloseParen,
			token.Whitespace,
			token.CloseParen,
		)
	})

	t.Run("function", func(t *testing.T) {
		CanTokenize(t, `
		( module 
			(func $add (param $lhs i32) (param $rhs i32) (result i32)
				local.get $lhs
				local.get $rhs
				i32.add))`,
			// ( module \n
			token.Whitespace,
			token.OpenParen,
			token.Whitespace,
			token.Reserved,
			token.Whitespace,
			// ( func $add
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.Id,
			token.Whitespace,
			// (param $lhs i32 )
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.Id,
			token.Whitespace,
			token.Reserved,
			token.CloseParen,
			token.Whitespace,
			// (param $rhs i32)
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.Id,
			token.Whitespace,
			token.Reserved,
			token.CloseParen,
			token.Whitespace,
			// (result i32)
			token.OpenParen,
			token.Reserved,
			token.Whitespace,
			token.Reserved,
			token.CloseParen,
			token.Whitespace,
			// local.get $lhs
			token.Reserved,
			token.Whitespace,
			token.Id,
			token.Whitespace,
			// local.get $rhs
			token.Reserved,
			token.Whitespace,
			token.Id,
			token.Whitespace,
			// i32.add))
			token.Reserved,
			token.CloseParen,
			token.CloseParen,
		)
	})
}

func CanTokenize(t *testing.T, input string, sequence ...token.Type) {
	lexer := lex.New2([]rune(input))
	for i, item := range sequence {
		tok, err := lexer.Next()
		require.Nil(t, err)
		require.Equal(t, item, tok.Type, "at %d", i)
	}
	tok, err := lexer.Next()
	require.Nil(t, err)
	require.Equal(t, token.EndOfStream, tok.Type, "at %d", len(sequence)-1)
}
