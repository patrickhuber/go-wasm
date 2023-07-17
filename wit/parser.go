// package wit covers parsing and generating wit files
// https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
package wit

import (
	"io"

	"github.com/patrickhuber/go-wasm/wit/ast"
)

func Parse(reader io.Reader) (ast.Ast, error) {
	lexer := NewLexer(reader)
	return parseAst(lexer)
}

func parseAst(lexer Lexer) (ast.Ast, error) {

	return &ast.AstNode{}, nil
}
