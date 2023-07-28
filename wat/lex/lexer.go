package lex

import (
	"fmt"
	"unicode"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wat/token"
)

type Lexer struct {
	input     []rune
	offset    int
	position  int
	column    int
	line      int
	peekToken *token.Token
}

func New(input []rune) *Lexer {
	return &Lexer{
		input: input,
	}
}

func (l *Lexer) Peek() (*token.Token, error) {
	// always return the peek token if it exists
	if l.peekToken != nil {
		return l.peekToken, nil
	}

	l.peekToken = l.next().Unwrap()
	return l.peekToken, nil
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
	case r == '(':
		if l.eat(';').Unwrap() {
			depth := 1
			for depth > 0 {
				r, ok := l.readRune().Deconstruct()
				if !ok {
					return result.Error[*token.Token](l.lexerError())
				}
				switch r {
				case '(':
					if l.eat(';').Unwrap() {
						depth++
					}
				case ';':
					if l.eat(')').Unwrap() {
						depth--
					}
				}
			}
			return l.token(token.BlockComment)
		}
		return l.token(token.OpenParen)
	case r == ')':
		return l.token(token.CloseParen)
	case r == '"':
		// does go handle the translation of escapes?
		for !l.eat('"').Unwrap() {
		}
		return l.token(token.String)
	case r == ';':
		if l.eat(';').Unwrap() {
			// line comment
			// ;; ... \n
			for l.eatIf(func(r rune) bool { return r != '\n' }).Unwrap() {
			}
			return l.token(token.LineComment)
		}
		return l.token(token.Reserved)
	case r == ',' || r == '[' || r == ']' || r == '{' || r == '}':
		return l.token(token.Reserved)
	case r == '$':
		for l.eatIf(isIdChar).Unwrap() {
		}
		return l.token(token.Id)
	case unicode.IsDigit(r):
		for l.eatIf(unicode.IsDigit).Unwrap() {

		}
		return l.token(token.Integer)
	case isIdChar(r) || r == '"':
		for l.eatIf(isIdChar).Unwrap() || l.eat('"').Unwrap() {
		}
		return l.token(token.Reserved)
	}

	return result.Errorf[*token.Token]("%w : unrecognized character '%c'", l.lexerError(), r)
}

func (l *Lexer) token(ty token.Type) types.Result[*token.Token] {

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
		}
	}

	// update the current offset to the position
	l.offset = l.position

	return result.Ok(tok)
}

var idCharMap = map[rune]struct{}{
	'!':  {},
	'#':  {},
	'$':  {},
	'%':  {},
	'&':  {},
	'\'': {},
	'*':  {},
	'+':  {},
	'-':  {},
	'.':  {},
	'/':  {},
	':':  {},
	'<':  {},
	'=':  {},
	'>':  {},
	'?':  {},
	'@':  {},
	'\\': {},
	'^':  {},
	'_':  {},
	'`':  {},
	'|':  {},
	'~':  {},
}

func isIdChar(ch rune) bool {
	_, ok := idCharMap[ch]
	if ok {
		return true
	}
	switch {
	case unicode.IsSpace(ch):
		return false
	case '0' <= ch && ch <= '9':
		return true
	case 'A' <= ch && ch <= 'Z':
		return true
	case 'a' <= ch && ch <= 'z':
		return true
	}
	return false
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

func (l *Lexer) expect(ch rune) (res types.Result[any]) {
	defer handle.Error(&res)

	r := l.readRune().Unwrap()
	if r != ch {
		return result.Errorf[any]("expected '%c' but found '%c'", ch, r)
	}
	return result.Ok[any](nil)
}

func (l *Lexer) expectIf(f func(ch rune) bool) (res types.Result[any]) {
	defer handle.Error(&res)

	r := l.readRune().Unwrap()
	if !f(r) {
		return result.Errorf[any]("expected '%c' but found '%c'", l.input[l.position-1], r)
	}
	return result.Ok[any](nil)
}

func (l *Lexer) peekRune() (op types.Option[rune]) {
	if l.position >= len(l.input) {
		return option.None[rune]()
	}
	r := l.input[l.position]
	return option.Some(r)
}

func (l *Lexer) readRune() types.Option[rune] {
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
