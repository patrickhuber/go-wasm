package parse

import (
	"fmt"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wat/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"

	"github.com/patrickhuber/go-wasm/wat/token"
)

func Parse(lexer *lex.Lexer) (ast.Ast, error) {
	return parse(lexer).Deconstruct()
}

func parse(lexer *lex.Lexer) (res types.Result[ast.Ast]) {
	defer handle.Error(&res)
	expect(lexer, token.OpenParen).Unwrap()
	tok := peek(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Errorf[ast.Ast]("%w : unrecognized token", parseError(tok))
	}
	var root ast.Ast
	switch tok.Capture {
	case "module":
		root = parseModule(lexer).Unwrap()
	case "component":
		root = parseComponent(lexer).Unwrap()
	default:
		return result.Errorf[ast.Ast](
			"%w : expected module, component but found '%s'",
			parseError(tok), tok.Capture)
	}
	expect(lexer, token.CloseParen).Unwrap()
	return result.Ok(root)
}

func parseModule(lexer *lex.Lexer) (res types.Result[*ast.Module]) {
	defer handle.Error(&res)
	expect(lexer, token.Reserved).Unwrap()
	return result.Ok(&ast.Module{})
}

func parseComponent(lexer *lex.Lexer) (res types.Result[*ast.Component]) {
	defer handle.Error(&res)
	expect(lexer, token.Reserved).Unwrap()
	return result.Ok(&ast.Component{})
}

func eat(lexer *lex.Lexer, ty token.Type) (res types.Result[bool]) {
	defer handle.Error(&res)

	tok := peek(lexer).Unwrap()
	if tok.Type != ty {
		return result.Ok(false)
	}

	expect(lexer, ty).Unwrap()
	return result.Ok(true)
}

func expect(lexer *lex.Lexer, ty token.Type) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type == ty {
		return result.Ok[any](nil)
	}
	return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), ty, tok.Type)
}

func next(lexer *lex.Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		res = result.New(lexer.Next())
		tok := res.Unwrap()
		switch tok.Type {
		// skip whitespace
		case token.Whitespace:
			continue
		// skip comments
		case token.BlockComment:
			continue
		// skip comments
		case token.LineComment:
			continue
		}
		return
	}
}

func peek(lexer *lex.Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		p := result.New(lexer.Peek())
		r := p.Unwrap()
		if r.Type != token.Whitespace {
			return p
		}
		// consume whitespace
		_ = result.New(lexer.Next()).Unwrap()
	}
}

func parseError(tok *token.Token) error {
	line := tok.Line + 1
	col := tok.Column + 1
	return fmt.Errorf(
		"error parsing at line %d, column %d, position %d",
		line,
		col,
		tok.Position)
}
