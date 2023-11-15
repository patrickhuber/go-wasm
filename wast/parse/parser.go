package parse

import (
	"fmt"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wast/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"
	watparse "github.com/patrickhuber/go-wasm/wat/parse"
	"github.com/patrickhuber/go-wasm/wat/token"
)

// Parse parses the wast spec https://github.com/WebAssembly/spec/tree/master/interpreter/#scripts
func Parse(input string) (*ast.Wast, error) {
	lexer := lex.New(input)
	return parseWast(lexer).Deconstruct()
}

func parseWast(lexer *lex.Lexer) (res types.Result[*ast.Wast]) {
	defer handle.Error(&res)
	var directives []ast.Directive
	for {
		tok := peek(lexer).Unwrap()
		if tok.Type == token.EndOfStream {
			break
		}
		directive := parseDirective(lexer).Unwrap()
		directives = append(directives, directive)
	}
	return result.Ok(&ast.Wast{
		Directives: directives,
	})
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
	case "module", "component":
		dir = ast.WatDirective{
			Wat: parseQuoteWat(lexer).Unwrap(),
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
	expectValue(lexer, token.Reserved, "assert_return").Unwrap()

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

	// assert_invalid
	expectValue(lexer, token.Reserved, "assert_invalid").Unwrap()

	module := parseQuoteWat(lexer).Unwrap()
	failure := parseString(lexer).Unwrap()

	return result.Ok[ast.Directive](ast.AssertInvalid{
		Module:  module,
		Failure: failure,
	})
}

func parseAssertMalformed(lexer *lex.Lexer) (res types.Result[ast.Directive]) {

	// (assert_malformed <module> <failure> )
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "assert_malformed").Unwrap()

	module := parseQuoteWat(lexer).Unwrap()
	failure := parseString(lexer).Unwrap()

	return result.Ok[ast.Directive](
		ast.AssertMalformed{
			Module:  module,
			Failure: failure,
		},
	)
}

func parseQuoteWat(lexer *lex.Lexer) (res types.Result[ast.QuoteWat]) {
	defer handle.Error(&res)

	// we need to look ahead 2 for the word 'quote'
	clone := lexer.Clone()

	expect(clone, token.OpenParen).Unwrap()

	var ty token.Type
	if eat(clone, token.Component).Unwrap() {
		ty = token.Component
	} else if eat(clone, token.Module).Unwrap() {
		ty = token.Module
	} else {
		tok := next(clone).Unwrap()
		return result.Errorf[ast.QuoteWat]("error parsing QuoteWat : %w", parseError(tok))
	}

	// 'quote'
	tok := next(clone).Unwrap()
	if tok.Type != token.Reserved || tok.Capture != "quote" {
		// this is a regular wat module, throw away the clone
		var wat ast.QuoteWat = parseWat(lexer).Unwrap()
		return result.Ok(wat)
	}

	// we are in a (module quote "") or (component quote "")
	// so we need to use the clone as the new lexer
	*lexer = *clone

	var quoteWat ast.QuoteWat
	switch ty {
	case token.Component:
		quoteWat = &ast.QuoteComponent{
			Quote: parseString(lexer).Unwrap(),
		}
	case token.Module:
		quoteWat = &ast.QuoteModule{
			Quote: parseString(lexer).Unwrap(),
		}
	default:
		return result.Errorf[ast.QuoteWat]("error parsing QuoteWat : unrecognized directive '%s'", ty)
	}

	expect(lexer, token.CloseParen).Unwrap()

	return result.Ok(quoteWat)

}

func parseWat(lexer *lex.Lexer) (res types.Result[*ast.Wat]) {
	wat, err := watparse.Parse(lexer)
	if err != nil {
		return result.Error[*ast.Wat](err)
	}
	return result.Ok(&ast.Wat{
		Wat: wat,
	})
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

	str := result.New(watparse.ParseString(lexer)).Unwrap()

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
		// some times constants can be ambiguous like '0'
		var f32 float32
		if peekFloat(lexer).Unwrap() {
			f32 = parseFloat32(lexer).Unwrap()
		} else {
			f32 = float32(parseInt32(lexer).Unwrap())
		}
		c = ast.F32Const{
			Value: f32,
		}
	case "f64.const":
		// some times constants can be ambiguous like '0'
		var f64 float64
		if peekFloat(lexer).Unwrap() {
			f64 = parseFloat64(lexer).Unwrap()
		} else if peekInteger(lexer).Unwrap() {
			f64 = float64(parseInt64(lexer).Unwrap())
		}
		c = ast.F64Const{
			Value: f64,
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

func peekFloat(lexer *lex.Lexer) (res types.Result[bool]) {
	tok, err := watparse.Peek(lexer)
	if err != nil {
		return result.New(false, err)
	}
	return result.Ok(tok.Type == token.Float)
}

func parseFloat32(lexer *lex.Lexer) (res types.Result[float32]) {
	return result.New(watparse.ParseFloat32(lexer))
}

func parseFloat64(lexer *lex.Lexer) (res types.Result[float64]) {
	return result.New(watparse.ParseFloat64(lexer))
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

func expectValue(lexer *lex.Lexer, ty token.Type, value string) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != ty {
		return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), ty, tok.Type)
	}
	if tok.Capture != value {
		return result.Errorf[any]("%w. expected '%v' value '%s' but found '%v' value '%s' ", parseError(tok), ty, value, tok, tok.Capture)
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
		if r.Type != token.Whitespace &&
			r.Type != token.BlockComment &&
			r.Type != token.LineComment {
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

func parseString(lexer *lex.Lexer) types.Result[string] {
	str, err := watparse.ParseString(lexer)
	return result.New(str, err)
}

func peekInteger(lexer *lex.Lexer) types.Result[bool] {
	i, err := watparse.Peek(lexer)
	if err != nil {
		return result.New(false, err)
	}
	return result.Ok(i.Type == token.Integer)
}

func parseInt32(lexer *lex.Lexer) types.Result[int32] {
	i, err := watparse.ParseInt32(lexer)
	return result.New(i, err)
}

func parseInt64(lexer *lex.Lexer) types.Result[int64] {
	i, err := watparse.ParseInt64(lexer)
	return result.New(i, err)
}
