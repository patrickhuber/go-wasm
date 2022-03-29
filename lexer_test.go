package wasm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-wasm"
)

var _ = Describe("Lexer", func() {
	It("can tokenize module", func() {
		CanTokenize("(module)",
			wasm.OpenParen,
			wasm.String,
			wasm.CloseParen)
	})
	It("can tokenize whitespace", func() {
		CanTokenize(" (module ) ",
			wasm.Whitespace,
			wasm.OpenParen,
			wasm.String,
			wasm.Whitespace,
			wasm.CloseParen,
			wasm.Whitespace)
	})
	It("can tokenize function", func() {
		CanTokenize(`
		( module 
			(func $add (param $lhs i32) (param $rhs i32) (result i32)
    			local.get $lhs
    			local.get $rhs
    			i32.add))`,
			wasm.Whitespace,
			wasm.OpenParen,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,

			wasm.OpenParen,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,

			wasm.OpenParen,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.CloseParen,
			wasm.Whitespace,

			wasm.OpenParen,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.CloseParen,
			wasm.Whitespace,

			wasm.OpenParen,
			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.CloseParen,
			wasm.Whitespace,

			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,

			wasm.String,
			wasm.Whitespace,
			wasm.String,
			wasm.Whitespace,

			wasm.String,
			wasm.CloseParen,

			wasm.CloseParen,
		)
	})
})

func CanTokenize(input string, sequence ...wasm.TokenType) {
	lexer := wasm.NewLexer(input)
	for i, item := range sequence {
		token, err := lexer.Next()
		Expect(err).To(BeNil())
		Expect(token.Type).To(Equal(item), "at %d", i)
	}
	token, err := lexer.Next()
	Expect(err).To(BeNil())
	Expect(token.Type).To(Equal(wasm.EndOfStream), "at %d", len(sequence)-1)
}
