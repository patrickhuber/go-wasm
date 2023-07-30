package parse

import (
	"fmt"
	"strconv"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
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

	m := &ast.Module{}
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := peek(lexer).Unwrap()
		switch tok.Capture {
		case "func":
			f := parseFunc(lexer).Unwrap()
			m.Functions = append(m.Functions, f)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok(m)
}

func parseComponent(lexer *lex.Lexer) (res types.Result[*ast.Component]) {
	defer handle.Error(&res)
	expect(lexer, token.Reserved).Unwrap()
	return result.Ok(&ast.Component{})
}

func parseFunc(lexer *lex.Lexer) (res types.Result[*ast.Function]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "func").Unwrap()

	function := &ast.Function{}
	if tok := peek(lexer).Unwrap(); tok.Type == token.Id {
		id := peek(lexer).Unwrap()
		if id.Type == token.Id {
			function.ID = option.Some(tok.Capture)
		} else {
			function.ID = option.None[string]()
		}
	}
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := peek(lexer).Unwrap()
		switch tok.Capture {
		case "param":
			param := parseParameter(lexer).Unwrap()
			function.Parameters = append(function.Parameters, *param)
		case "result":
			result := parseResult(lexer).Unwrap()
			function.Results = append(function.Results, *result)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}
	for {
		tok := peek(lexer).Unwrap()
		if tok.Type == token.CloseParen {
			break
		}
		instruction := parseInstruction(lexer).Unwrap()
		function.Instructions = append(function.Instructions, instruction)
	}
	return result.Ok(function)
}

func parseParameter(lexer *lex.Lexer) (res types.Result[*ast.Parameter]) {
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "param").Unwrap()

	parameter := &ast.Parameter{}
	tok := peek(lexer).Unwrap()
	if tok.Type == token.Id {
		id := parseId(lexer).Unwrap()
		parameter.ID = option.Some(id)
	} else {
		parameter.ID = option.None[string]()
	}
	parameter.Type = parseType(lexer).Unwrap()
	return result.Ok(parameter)
}

func parseResult(lexer *lex.Lexer) (res types.Result[*ast.Result]) {
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "result").Unwrap()

	return result.Ok(&ast.Result{
		Type: parseType(lexer).Unwrap(),
	})
}

func parseInstruction(lexer *lex.Lexer) (res types.Result[ast.Instruction]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Error[ast.Instruction](parseError(tok))
	}
	var inst ast.Instruction
	switch tok.Capture {
	case "local.get":
		inst = ast.LocalGet{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "i32.add":
		inst = ast.I32Add{}
	default:
		return result.Errorf[ast.Instruction]("%w : error parsing instruction. Unrecognized instruction %v : %s", parseError(tok), tok.Type, tok.Capture)
	}
	return result.Ok(inst)
}

func parseIndex(lexer *lex.Lexer) (res types.Result[ast.Index]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	var index ast.Index
	switch tok.Type {
	case token.Integer:
		i := result.New(strconv.Atoi(tok.Capture)).Unwrap()
		index = &ast.RawIndex{
			Index: uint32(i),
		}
	case token.Id:
		index = &ast.IDIndex{
			ID: tok.Capture,
		}
	default:
		return result.Errorf[ast.Index]("%w : error parsing index", parseError(tok))
	}
	return result.Ok(index)
}

func parseType(lexer *lex.Lexer) (res types.Result[ast.Type]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	var ty ast.Type
	// todo: enhance the lexer to parse these as tokens
	switch tok.Capture {
	case "i32":
		ty = ast.I32{}
	case "i64":
		ty = ast.I64{}
	case "f32":
		ty = ast.F32{}
	case "f64":
		ty = ast.F64{}
	default:
		return result.Errorf[ast.Type]("%w : error parsing type. expected (i32, i64, f32, f64) but found %s", parseError(tok), tok.Capture)
	}
	return result.Ok(ty)
}

func parseId(lexer *lex.Lexer) (res types.Result[string]) {
	tok := next(lexer).Unwrap()
	if tok.Type != token.Id {
		return result.Errorf[string]("%w", parseError(tok))
	}
	return result.Ok(tok.Capture)
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

func expectValue(lexer *lex.Lexer, ty token.Type, capture string) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != ty {
		return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), ty, tok.Type)
	}
	if tok.Capture != capture {
		return result.Errorf[any]("%w. expected type:'%v' value:'%s' but found type:'%v' value:'%s'",
			parseError(tok), ty, capture, tok.Type, tok.Capture)
	}
	return result.Ok[any](nil)
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
