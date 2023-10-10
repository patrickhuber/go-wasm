package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wast/ast"
	wat "github.com/patrickhuber/go-wasm/wat/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"
	watparse "github.com/patrickhuber/go-wasm/wat/parse"
	"github.com/patrickhuber/go-wasm/wat/token"
)

// Parse parses the wast spec https://github.com/WebAssembly/spec/tree/master/interpreter/#scripts
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
		*lexer = *clone
		dir = parseAssertReturn(lexer).Unwrap()
	case "assert_invalid":
		*lexer = *clone
		dir = parseAssertInvalid(lexer).Unwrap()
	case "assert_malformed":
		*lexer = *clone
		dir = parseAssertMalformed(lexer).Unwrap()
	case "assert_trap":
		*lexer = *clone
		dir = parseAssertTrap(lexer).Unwrap()
	default:
		return result.Error[ast.Directive](parseError(tok))
	}

	expect(lexer, token.CloseParen).Unwrap()
	return result.Ok(dir)
}

func parseAssertReturn(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	defer handle.Error(&res)

	// ( assert_return <action> <result>* )

	// assert_return
	expect(lexer, token.Reserved).Unwrap()

	// <action>
	action := parseAction(lexer).Unwrap()
	// <result>*
	results := parseResults(lexer).Unwrap()

	return result.Ok[ast.Directive](ast.AssertReturn{
		Action:  action,
		Results: results,
	})
}

func parseAssertInvalid(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	// ( assert_invalid <module> <failure> )
	defer handle.Error(&res)

	// ( assert_invalid <module> <failure> )

	// assert_invalid
	expect(lexer, token.Reserved).Unwrap()

	m, err := watparse.Parse(lexer)
	if err != nil {
		return result.Error[ast.Directive](err)
	}
	module, ok := m.(*wat.Module)
	if !ok {
		return result.Errorf[ast.Directive]("expected module but found %T", m)
	}

	failure := parseString(lexer).Unwrap()

	return result.Ok[ast.Directive](ast.AssertInvalid{
		Module:  module,
		Failure: failure,
	})
}

func parseAssertMalformed(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	return result.Errorf[ast.Directive]("parseAssertMalformed not implemented")
}

func parseAssertTrap(lexer *lex.Lexer) (res types.Result[ast.Directive]) {
	defer handle.Error(&res)

	// assert_return
	expect(lexer, token.Reserved).Unwrap()

	// ( assert_trap <module> <failure> )
	action := parseAction(lexer).Unwrap()
	failure := parseString(lexer).Unwrap()

	return result.Ok[ast.Directive](ast.AssertTrap{
		Action:  action,
		Failure: failure,
	})
}

func parseAction(lexer *lex.Lexer) (res types.Result[ast.Action]) {
	defer handle.Error(&res)

	/*	action:
		( invoke <name>? <string> <const>* )       ;; invoke function export
		( get <name>? <string> )                   ;; get global export
	*/
	expect(lexer, token.OpenParen).Unwrap()

	tok := next(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Errorf[ast.Action]("%w : unrecognized token", parseError(tok))
	}

	var action ast.Action
	switch tok.Capture {
	case "invoke":
		action = parseInvoke(lexer).Unwrap()
	case "get":
		action = parseGet(lexer).Unwrap()
	}

	expect(lexer, token.CloseParen).Unwrap()

	return result.Ok(action)
}

func parseInvoke(lexer *lex.Lexer) (res types.Result[ast.Invoke]) {
	defer handle.Error(&res)

	tok := peek(lexer).Unwrap()

	var name types.Option[string]
	if tok.Type == token.Id {
		name = option.Some(tok.Capture)
		expect(lexer, token.Id).Unwrap()
	} else {
		name = option.None[string]()
	}

	str := parseString(lexer).Unwrap()

	consts := parseConsts(lexer).Unwrap()

	return result.Ok(ast.Invoke{
		Name:   name,
		String: str,
		Const:  consts,
	})
}

func parseConsts(lexer *lex.Lexer) (res types.Result[[]ast.Const]) {
	defer handle.Error(&res)
	var consts []ast.Const
	for eat(lexer, token.OpenParen).Unwrap() {
		c := parseConst(lexer).Unwrap()
		consts = append(consts, c)

		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok(consts)
}

func parseConst(lexer *lex.Lexer) (res types.Result[ast.Const]) {
	/*
		const:
		  ( <num_type>.const <num> )                 ;; number value
		  ( <vec_type> <vec_shape> <num>+ )          ;; vector value
		  ( ref.null <ref_kind> )                    ;; null reference
		  ( ref.extern <nat> )                       ;; host reference
	*/
	tok := next(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Error[ast.Const](parseError(tok))
	}

	var c ast.Const
	switch tok.Capture {
	case "i32.const":
		c = ast.I32Const{
			Value: parseInt32(lexer).Unwrap(),
		}
	case "i64.const":
		c = ast.I64Const{
			Value: parseInt64(lexer).Unwrap(),
		}
	case "f32.const":
		c = ast.F32Const{
			Value: parseFloat32(lexer).Unwrap(),
		}
	case "f64.const":
		c = ast.F64Const{
			Value: parseFloat64(lexer).Unwrap(),
		}
	case "ref.null":
	case "ref.extern":
	default:
		return result.Error[ast.Const](parseError(tok))
	}

	return result.Ok(c)
}

func parseGet(lexer *lex.Lexer) (res types.Result[ast.Get]) {
	return result.Errorf[ast.Get]("parseGet not implemented")
}

func parseResults(lexer *lex.Lexer) (res types.Result[[]ast.Result]) {
	defer handle.Error(&res)

	var results []ast.Result
	for eat(lexer, token.OpenParen).Unwrap() {
		c := parseConst(lexer).Unwrap()
		r, ok := c.(ast.Result)
		if !ok {
			continue
		}
		results = append(results, r)
		expect(lexer, token.CloseParen).Unwrap()
	}

	return result.Ok(results)
}

func parseString(lexer *lex.Lexer) (res types.Result[string]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	if tok.Type != token.String {
		return result.Errorf[string]("%w expected string found %s", parseError(tok), tok.Type)
	}
	return result.Ok(strings.Trim(tok.Capture, "\""))
}

func parseInt32(lexer *lex.Lexer) (res types.Result[int32]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	if tok.Type != token.Integer {
		return result.Errorf[int32]("%w expected integer found %s", parseError(tok), tok.Type)
	}

	var v int32
	// 0x80000000 will overflow if we don't parse it as uint32 so
	// any negative number will be parsed as a regular int and everything
	// else is parsed as int
	if strings.HasPrefix(tok.Capture, "-") {
		i, err := strconv.ParseInt(tok.Capture, 0, 32)
		if err != nil {
			return result.Errorf[int32]("%w", err)
		}
		v = int32(i)
	} else {
		u, err := strconv.ParseUint(tok.Capture, 0, 32)
		if err != nil {
			return result.Errorf[int32]("%w", err)
		}
		v = int32(u)
	}
	return result.Ok(v)
}

func parseInt64(lexer *lex.Lexer) (res types.Result[int64]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	if tok.Type != token.Integer {
		return result.Errorf[int64]("%w", parseError(tok))
	}

	var v int64
	// 0x80000000_00000000 will overflow if we don't parse it as uint64 so
	// any negative number will be parsed as a regular int and everything
	// else is parsed as int
	if strings.HasPrefix(tok.Capture, "-") {
		f, err := strconv.ParseInt(tok.Capture, 0, 64)
		if err != nil {
			return result.Errorf[int64]("%w", err)
		}
		v = int64(f)
	} else {
		f, err := strconv.ParseUint(tok.Capture, 0, 64)
		if err != nil {
			return result.Errorf[int64]("%w", err)
		}
		v = int64(f)
	}
	return result.Ok(v)
}

func parseFloat32(lexer *lex.Lexer) (res types.Result[float32]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	if tok.Type != token.Float {
		return result.Errorf[float32]("%w", parseError(tok))
	}
	f, err := strconv.ParseFloat(tok.Capture, 32)
	if err != nil {
		return result.Errorf[float32]("%w", err)
	}
	return result.Ok(float32(f))
}

func parseFloat64(lexer *lex.Lexer) (res types.Result[float64]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	if tok.Type != token.Float {
		return result.Errorf[float64]("%w", parseError(tok))
	}
	f, err := strconv.ParseFloat(tok.Capture, 64)
	if err != nil {
		return result.Errorf[float64]("%w", err)
	}
	return result.Ok(f)
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
