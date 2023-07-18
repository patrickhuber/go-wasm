// package wit covers parsing and generating wit files
// https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
package wit

import (
	"fmt"
	"io"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	abi "github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/wit/ast"
	"github.com/patrickhuber/go-wasm/wit/token"
)

func Parse(reader io.Reader) (*ast.Ast, error) {
	lexer := NewLexer(reader)
	return parseAst(lexer).Deconstruct()
}

func parseAst(lexer Lexer) (res types.Result[*ast.Ast]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	n := &ast.Ast{}

	switch tok.Capture {
	case "package":
		packageName := parsePackageName(lexer).Unwrap()
		n.PackageName = option.Some(*packageName)
		tok = next(lexer).Unwrap()
	default:
		n.PackageName = option.None[ast.PackageName]()
	}

	for {
		item := &ast.AstItem{}
		switch tok.Capture {
		case "use":
			item.Use = parseTopLevelUse(lexer).Unwrap()
		case "world":
			item.World = parseWorld(lexer).Unwrap()
		case "interface":
			item.Interface = parseInterface(lexer).Unwrap()
		}
		n.Items = append(n.Items, *item)

		tok = next(lexer).Unwrap()
		if tok.Type == token.EndOfStream {
			break
		}
	}
	return result.Ok(n)
}

func parsePackageName(lexer Lexer) (res types.Result[*ast.PackageName]) {
	defer handle.Error(&res)

	// id
	tok := next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()

	packageName := &ast.PackageName{
		Namespace: tok.Capture,
	}

	// ':'
	tok = next(lexer).Unwrap()
	expect(tok, token.Colon).Unwrap()

	// id
	tok = next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()
	packageName.Name = tok.Capture

	peek := peek(lexer).Unwrap()
	if !optional(peek, token.At) {
		packageName.Version = option.None[ast.Version]()
		return result.Ok(packageName)
	}

	tok = next(lexer).Unwrap()
	expect(tok, token.At).Unwrap()

	version := parseVersion(lexer).Unwrap()
	packageName.Version = option.Some(*version)

	return result.Ok(packageName)
}

func parseVersion(lexer Lexer) (res types.Result[*ast.Version]) {
	defer handle.Error(&res)
	version := &ast.Version{}
	return result.Ok(version)
}

func parseTopLevelUse(lexer Lexer) (res types.Result[*ast.TopLevelUse]) {
	defer handle.Error(&res)
	topLevelUse := &ast.TopLevelUse{}
	return result.Ok(topLevelUse)
}

func parseInterface(lexer Lexer) (res types.Result[*ast.Interface]) {

	defer handle.Error(&res)

	// id
	tok := next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()

	// '{'
	tok = next(lexer).Unwrap()
	expect(tok, token.OpenBrace).Unwrap()

	inter := &ast.Interface{
		Name: tok.Capture,
	}

	for {
		inter.Items = append(inter.Items, *parseInterfaceItem(lexer).Unwrap())

		peekTok := peek(lexer).Unwrap()

		// exit on '}'
		if optional(peekTok, token.CloseBrace) {
			break
		}
	}

	tok = next(lexer).Unwrap()

	// '}'
	expect(tok, token.CloseBrace).Unwrap()

	return result.Ok(inter)
}

func parseInterfaceItem(lexer Lexer) (res types.Result[*ast.InterfaceItem]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()

	item := &ast.InterfaceItem{}

	switch tok.Capture {
	case "use":
	case "resource":
	case "variant":
	case "record":
	case "union":
	case "flags":
	case "enum":
	case "type":

	default:
		// tok == id
		funcItem := parseNamedFunc(tok, lexer).Unwrap()
		item.Func = funcItem
	}
	return result.Ok(item)
}

func parseNamedFunc(name *token.Token, lexer Lexer) (res types.Result[*ast.NamedFunc]) {

	defer handle.Error(&res)

	named := &ast.NamedFunc{
		Name: name.Capture,
	}

	tok := next(lexer).Unwrap()
	expect(tok, token.Colon).Unwrap()

	_func := parseFunc(lexer).Unwrap()
	named.Func = _func

	return result.Ok(named)
}

func parseFunc(lexer Lexer) (res types.Result[*ast.Func]) {

	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	expectValue(tok, token.String, "func").Unwrap()

	tok = next(lexer).Unwrap()
	expect(tok, token.OpenParen).Unwrap()

	_func := &ast.Func{}
	for {
		// name ':' type
		param := parseParameter(lexer).Unwrap()
		_func.Params = append(_func.Params, *param)

		peekTok := peek(lexer).Unwrap()

		if peekTok.Type == token.Comma {
			tok = next(lexer).Unwrap()
			expect(tok, token.Comma).Unwrap()
		} else if peekTok.Type == token.CloseParen {
			break
		} else {
			return result.Errorf[*ast.Func]("%w. expected ',' or ')' but found %s", parseError(peekTok), peekTok.Capture)
		}

		peekTok = peek(lexer).Unwrap()
		if peekTok.Type == token.CloseParen {
			break
		}
	}

	tok = next(lexer).Unwrap()
	expect(tok, token.CloseParen).Unwrap()

	return result.Ok(_func)
}

func parseParameter(lexer Lexer) (res types.Result[*ast.Parameter]) {
	defer handle.Error(&res)

	parameter := &ast.Parameter{}

	tok := next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()

	parameter.Id = tok.Capture

	tok = next(lexer).Unwrap()
	expect(tok, token.Colon).Unwrap()

	parameter.Type = parseType(lexer).Unwrap()
	return result.Ok(parameter)
}

func parseType(lexer Lexer) (res types.Result[abi.Type]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	expect(tok, token.String).Unwrap()

	var ty abi.Type
	switch tok.Capture {
	case "string":
		ty = abi.NewString()
	default:
		return result.Errorf[abi.Type]("unrecognized type %s", tok.Capture)
	}
	return result.Ok(ty)
}

func parseWorld(lexer Lexer) (res types.Result[*ast.World]) {
	defer handle.Error(&res)
	return result.Errorf[*ast.World]("not implemented")
}

func next(lexer Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		res = result.New(lexer.Next())
		tok := res.Unwrap()
		if tok.Type != token.Whitespace {
			return
		}
	}
}

func peek(lexer Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		res = result.New(lexer.Next())
		tok := res.Unwrap()
		if tok.Type != token.Whitespace {
			return
		}
	}
}

func optional(tok *token.Token, tokenType token.TokenType) bool {
	return tok.Type == tokenType
}

func optionalValue(tok *token.Token, tokenType token.TokenType, capture string) bool {
	return optional(tok, tokenType) && tok.Capture == capture
}

func expect(tok *token.Token, tokenType token.TokenType) types.Result[any] {
	if optional(tok, tokenType) {
		return result.Ok[any](nil)
	}
	return result.Errorf[any]("%w. expected %v but found %v", parseError(tok), tokenType, tok.Type)
}

func expectValue(tok *token.Token, tokenType token.TokenType, capture string) types.Result[any] {
	if optionalValue(tok, tokenType, capture) {
		return result.Ok[any](nil)
	}
	return result.Errorf[any]("%w. expected %v but found %v", parseError(tok), tokenType, tok.Type)
}

func parseError(tok *token.Token) error {
	return fmt.Errorf("error parsing at position %d line %d column %d", tok.Position, tok.Line, tok.Column)
}
