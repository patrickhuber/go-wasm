package lex

import (
	"fmt"
	"unicode"

	"github.com/patrickhuber/go-wasm/wit/token"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
)

type Lexer struct {
	input     []rune
	offset    int
	position  int
	column    int
	line      int
	peekToken *token.Token
}

func (l *Lexer) Line() int {
	return l.line
}

func (l *Lexer) Column() int {
	return l.column
}

func New(input []rune) *Lexer {
	return &Lexer{
		input:    input,
		position: 0,
		column:   0,
		line:     0,
	}
}

func (l *Lexer) Next() (*token.Token, error) {
	// any peek token?
	if l.peekToken == nil {
		return l.next().Deconstruct()
	}

	tok := l.peekToken
	l.peekToken = nil

	return tok, nil
}

func (l *Lexer) Peek() (*token.Token, error) {

	// always return the peek token if it exists
	if l.peekToken != nil {
		return l.peekToken, nil
	}

	l.peekToken = l.next().Unwrap()
	return l.peekToken, nil
}

func (l *Lexer) next() (res types.Result[*token.Token]) {
	defer handle.Error(&res)

	r, ok := l.readRune().Deconstruct()
	if !ok {
		return l.token(token.EndOfStream)
	}

	switch {
	case unicode.IsSpace(r):
		for l.eatIf(unicode.IsSpace).Unwrap() {
		}
		return l.token(token.Whitespace)
	case r == '/':

		if l.eat('/').Unwrap() {
			// line comment
			for l.eatIf(func(r rune) bool { return r != '\n' }).Unwrap() {
			}
			return l.token(token.LineComment)
		} else if l.eat('*').Unwrap() {

			// else block comment
			depth := 1
			for depth > 0 {
				r, ok := l.readRune().Deconstruct()
				if !ok {
					return result.Error[*token.Token](l.lexerError())
				}
				switch r {
				case '/':
					if l.eat('*').Unwrap() {
						depth++
					}
				case '*':
					if l.eat('/').Unwrap() {
						depth--
					}
				}
			}
			return l.token(token.BlockComment)
		}
		return l.token(token.Slash)
	case r == '=':
		return l.token(token.Equal)
	case r == ',':
		return l.token(token.Comma)
	case r == ':':
		return l.token(token.Colon)
	case r == '.':
		return l.token(token.Period)
	case r == ';':
		return l.token(token.Semicolon)
	case r == '(':
		return l.token(token.OpenParen)
	case r == ')':
		return l.token(token.CloseParen)
	case r == '{':
		return l.token(token.OpenBrace)
	case r == '}':
		return l.token(token.CloseBrace)
	case r == '<':
		return l.token(token.Less)
	case r == '>':
		return l.token(token.Greater)
	case r == '*':
		return l.token(token.Star)
	case r == '@':
		return l.token(token.At)
	case r == '-':
		if l.eat('>').Unwrap() {
			return l.token(token.RightArrow)
		} else {
			return l.token(token.Minus)
		}
	case r == '+':
		return l.token(token.Plus)
	case r == '%':
		if l.eatIf(isKeyLikeStart).Unwrap() {
			for l.eatIf(isKeyLikeContinue).Unwrap() {
			}
		}
		return l.token(token.ExplicitId)

	case isKeyLikeStart(r):
		// identifier | string
		for l.eatIf(isKeyLikeContinue).Unwrap() {
		}
		runes := l.capture()
		ty, ok := keywordMap[string(runes)]
		if !ok {
			ty = token.Id
		}
		return l.token(ty)

	case unicode.IsDigit(r):
		for l.eatIf(unicode.IsDigit).Unwrap() {

		}
		return l.token(token.Integer)
	}

	return result.Errorf[*token.Token]("%w : unrecognized character %c", l.lexerError(), r)
}

var keywordMap = map[string]token.TokenType{
	"use":         token.Use,
	"type":        token.Type,
	"func":        token.Func,
	"u8":          token.U8,
	"u16":         token.U16,
	"u32":         token.U32,
	"u64":         token.U64,
	"s8":          token.S8,
	"s16":         token.S16,
	"s32":         token.S32,
	"s64":         token.S64,
	"float32":     token.Float32,
	"float64":     token.Float64,
	"char":        token.Char,
	"resource":    token.Resource,
	"own":         token.Own,
	"borrow":      token.Borrow,
	"record":      token.Record,
	"flags":       token.Flags,
	"variant":     token.Variant,
	"enum":        token.Enum,
	"union":       token.Union,
	"bool":        token.Bool,
	"string":      token.String,
	"option":      token.Option,
	"result":      token.Result,
	"future":      token.Future,
	"stream":      token.Stream,
	"list":        token.List,
	"_":           token.Underscore,
	"as":          token.As,
	"from":        token.From,
	"static":      token.Static,
	"interface":   token.Interface,
	"tuple":       token.Tuple,
	"world":       token.World,
	"import":      token.Import,
	"export":      token.Export,
	"package":     token.Package,
	"constructor": token.Constructor,
	"include":     token.Include,
	"with":        token.With,
}

func (l *Lexer) token(ty token.TokenType) types.Result[*token.Token] {

	// snapshot the state for the current token
	tok := &token.Token{
		Type:     ty,
		Position: l.offset,
		Column:   l.column,
		Line:     l.line,
		Runes:    l.input[l.offset:l.position],
	}

	// fast forward updating metrics
	for i := l.offset; i < l.position; i++ {
		ch := l.input[i]
		if ch == '\n' {
			l.line++
			l.column = 0
		} else {
			l.column++
		}
	}

	// update the current offset to the position
	l.offset = l.position

	return result.Ok(tok)
}

func (l *Lexer) capture() []rune {
	return l.input[l.offset:l.position]
}

func (l *Lexer) eatIf(f func(ch rune) bool) (res types.Result[bool]) {
	defer handle.Error(&res)

	p, ok := l.peekRune().Deconstruct()
	if !ok {
		return result.Ok(false)
	}
	if !f(p) {
		return result.Ok(false)
	}
	l.expectIf(f).Unwrap()
	return result.Ok(true)
}

func (l *Lexer) eat(ch rune) (res types.Result[bool]) {
	defer handle.Error(&res)

	p, ok := l.peekRune().Deconstruct()
	if !ok {
		return result.Ok(false)
	}
	if p != ch {
		return result.Ok(false)
	}
	l.expect(ch).Unwrap()
	return result.Ok(true)
}

func (l *Lexer) expectIf(f func(ch rune) bool) (res types.Result[any]) {
	defer handle.Error(&res)

	r := l.readRune().Unwrap()
	if !f(r) {
		return result.Errorf[any]("expected '%c' but found '%c'", l.input[l.position-1], r)
	}
	return result.Ok[any](nil)
}

func (l *Lexer) expect(ch rune) (res types.Result[any]) {
	defer handle.Error(&res)

	r := l.readRune().Unwrap()
	if r != ch {
		return result.Errorf[any]("expected '%c' but found '%c'", ch, r)
	}
	return result.Ok[any](nil)
}

func (l *Lexer) peekRune() (res types.Option[rune]) {
	if l.position >= len(l.input) {
		return option.None[rune]()
	}
	r := l.input[l.position]
	return option.Some(r)
}

func (l *Lexer) readRune() (res types.Option[rune]) {
	if l.position >= len(l.input) {
		return option.None[rune]()
	}
	r := l.input[l.position]
	l.position++
	return option.Some(r)
}

func (l *Lexer) lexerError() error {
	return fmt.Errorf("error parsing at line: %d column: %d position: %d", l.line, l.column, l.position)
}

func isKeyLikeStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || ch == '-'
}

func isKeyLikeContinue(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '-'
}
