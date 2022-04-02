package wat_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/go-wasm/wat"
)

var _ = Describe("Lexer", func() {
	It("can tokenize module", func() {
		CanTokenize("(module)",
			wat.OpenParen,
			wat.String,
			wat.CloseParen)
	})
	It("can tokenize whitespace", func() {
		CanTokenize(" (module ) ",
			wat.Whitespace,
			wat.OpenParen,
			wat.String,
			wat.Whitespace,
			wat.CloseParen,
			wat.Whitespace)
	})
	It("can tokenize function", func() {
		CanTokenize(`
		( module 
			(func $add (param $lhs i32) (param $rhs i32) (result i32)
    			local.get $lhs
    			local.get $rhs
    			i32.add))`,
			wat.Whitespace,
			wat.OpenParen,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,

			wat.OpenParen,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,

			wat.OpenParen,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.CloseParen,
			wat.Whitespace,

			wat.OpenParen,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.CloseParen,
			wat.Whitespace,

			wat.OpenParen,
			wat.String,
			wat.Whitespace,
			wat.String,
			wat.CloseParen,
			wat.Whitespace,

			wat.String,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,

			wat.String,
			wat.Whitespace,
			wat.String,
			wat.Whitespace,

			wat.String,
			wat.CloseParen,

			wat.CloseParen,
		)
	})
})

func CanTokenize(input string, sequence ...wat.TokenType) {
	lexer := wat.NewLexer(input)
	for i, item := range sequence {
		token, err := lexer.Next()
		Expect(err).To(BeNil())
		Expect(token.Type).To(Equal(item), "at %d", i)
	}
	token, err := lexer.Next()
	Expect(err).To(BeNil())
	Expect(token.Type).To(Equal(wat.EndOfStream), "at %d", len(sequence)-1)
}
