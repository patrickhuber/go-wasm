package parse

import (
	"fmt"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wast/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"
	watparse "github.com/patrickhuber/go-wasm/wat/parse"
	"github.com/patrickhuber/go-wasm/wat/token"
)

func Parse(input string) ([]ast.Directive, error) {
	return parse(input).Deconstruct()
}

func parse(input string) (res types.Result[[]ast.Directive]) {
	defer handle.Error(&res)

	var directives []ast.Directive
	lexer := lex.New(input)
	for {
		peek := result.New(lexer.Peek()).Unwrap()
		if peek.Type == token.EndOfStream {
			break
		}
		directive := parseDirective(lexer).Unwrap()
		directives = append(directives, directive)
	}

	return result.Ok(directives)
}

func parseDirective(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	defer handle.Error(&res)

	// create a clone of the lexer to check if this is a wat directive or a wast directive
	// revert to the original lexer if wat
	// use the clone if wast
	clone := lexer.Clone()

	var dir ast.Directive
	expect(clone, token.OpenParen).Unwrap()

	tok := peek(clone).Unwrap()

	switch tok.Capture {
	case "module":
		fallthrough
	case "component":
		dir = ast.Wat{
			Wat: result.New(watparse.Parse(lexer)).Unwrap(),
		}
		// exit early as wat parse will eat the last close paren
		return result.Ok(dir)
	case "assert_return":
		lexer = clone
		dir = parseAssertReturn(lexer).Unwrap()
	case "assert_invalid":
		lexer = clone
		dir = parseAssertInvalid(lexer).Unwrap()
	case "assert_malformed":
		lexer = clone
		dir = parseAssertMalformed(lexer).Unwrap()
	case "assert_trap":
		lexer = clone
		dir = parseAssertTrap(lexer).Unwrap()
	default:
		return result.Error[ast.Directive](parseError(tok))
	}

	expect(lexer, token.CloseParen).Unwrap()
	return result.Ok(dir)
}

func parseAssertReturn(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	return result.Errorf[ast.Directive]("not implemented")
}

func parseAssertInvalid(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	return result.Errorf[ast.Directive]("not implemented")
}

func parseAssertMalformed(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	return result.Errorf[ast.Directive]("not implemented")
}

func parseAssertTrap(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	return result.Errorf[ast.Directive]("not implemented")
}

func parseInvoke(lexer *lex.Lexer) (res types.Result[ast.Invoke]) {
	return result.Errorf[ast.Invoke]("not implemented")
}

type Soda interface {
	soda()
}
type soda int

func (soda) soda() {}

const (
	Fizz soda = 0
	Buzz soda = 1
	Baz  soda = 2
)

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
