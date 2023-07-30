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

	t.Run("wast_i32", func(t *testing.T) {
		content := `
;; i32 operations

(module
	(func (export "add") (param $x i32) (param $y i32) (result i32) (i32.add (local.get $x) (local.get $y)))
)
(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`
		CanTokenize(t, string(content),
			token.Whitespace,
			// ;; i32 operations
			token.LineComment,
			token.Whitespace,
			// (module
			token.OpenParen, token.Reserved, token.Whitespace,
			// (func
			token.OpenParen, token.Reserved, token.Whitespace,
			// (export "add")
			token.OpenParen, token.Reserved, token.Whitespace, token.String, token.CloseParen, token.Whitespace,
			// (param $x i32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (param $y i32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (result i32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (i32.add (local.get $x)
			token.OpenParen, token.Reserved, token.Whitespace, token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.CloseParen, token.Whitespace,
			// (local.get $y)))
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.CloseParen, token.CloseParen, token.CloseParen, token.Whitespace,
			// )
			token.CloseParen, token.Whitespace,
			// (assert_return
			token.OpenParen, token.Reserved, token.Whitespace,
			// (invoke "add"
			token.OpenParen, token.Reserved, token.Whitespace, token.String, token.Whitespace,
			// (i32.const 1)
			token.OpenParen, token.Reserved, token.Whitespace, token.Integer, token.CloseParen, token.Whitespace,
			// (i32.const 1))
			token.OpenParen, token.Reserved, token.Whitespace, token.Integer, token.CloseParen, token.CloseParen, token.Whitespace,
			// (i32.contt 2))
			token.OpenParen, token.Reserved, token.Whitespace, token.Integer, token.CloseParen, token.CloseParen,
		)
	})

	t.Run("wast_f32", func(t *testing.T) {
		content := `
;; Test all the f32 operators on major boundary values and all special
;; values (except comparison and bitwise operators, which are tested in
;; f32_bitwise.wast and f32_cmp.wast).

(module
	(func (export "add") (param $x f32) (param $y f32) (result f32) (f32.add (local.get $x) (local.get $y)))
)
(assert_return (invoke "add" (f32.const 0x0p+0) (f32.const -0x0p+0)) (f32.const 0x0p+0))
(assert_return (invoke "add" (f32.const -0x0p+0) (f32.const -inf)) (f32.const -inf))
(assert_return (invoke "add" (f32.const -0x0p+0) (f32.const -nan)) (f32.const nan:canonical))`
		CanTokenize(t, string(content),
			token.Whitespace,
			// ;; Test all the f32 operators on major boundary values and all special
			token.LineComment,
			// ;; values (except comparison and bitwise operators, which are tested in
			token.LineComment,
			// ;; f32_bitwise.wast and f32_cmp.wast).
			token.LineComment, token.Whitespace,
			// (module
			token.OpenParen, token.Reserved, token.Whitespace,
			// (func
			token.OpenParen, token.Reserved, token.Whitespace,
			// (export "add")
			token.OpenParen, token.Reserved, token.Whitespace, token.String, token.CloseParen, token.Whitespace,
			// (param $x f32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (param $y f32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (result i32)
			token.OpenParen, token.Reserved, token.Whitespace, token.Reserved, token.CloseParen, token.Whitespace,
			// (f32.add
			token.OpenParen, token.Reserved, token.Whitespace,
			// (local.get $x)
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.CloseParen, token.Whitespace,
			// (local.get $y)))
			token.OpenParen, token.Reserved, token.Whitespace, token.Id, token.CloseParen, token.CloseParen, token.CloseParen, token.Whitespace,
		)
	})
}

func CanTokenize(t *testing.T, input string, sequence ...token.Type) {
	lexer := lex.New(input)
	for i, item := range sequence {
		tok, err := lexer.Next()
		require.Nil(t, err)
		require.Equal(t, item, tok.Type, "at %d", i)
	}
	tok, err := lexer.Next()
	require.Nil(t, err)
	require.Equal(t, token.EndOfStream, tok.Type, "at %d", len(sequence)-1)
}
