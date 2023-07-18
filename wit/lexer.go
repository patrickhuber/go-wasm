package wit

import (
	"bufio"
	"io"
	"strings"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wit/token"
)

type Lexer interface {
	Next() (*token.Token, error)
	Peek() (*token.Token, error)
}

type lexer struct {
	reader   *bufio.Reader
	peekTok  *token.Token
	position int
	line     int
	column   int
}

func NewLexer(reader io.Reader) Lexer {
	return &lexer{
		reader:   bufio.NewReader(reader),
		peekTok:  nil,
		position: 0,
		line:     0,
		column:   0,
	}
}

func (l *lexer) Next() (*token.Token, error) {
	return l.next().Deconstruct()
}

func (l *lexer) next() (res types.Result[*token.Token]) {
	defer handle.Error(&res)

	if l.peekTok != nil {
		ret := l.peekTok
		l.peekTok = nil
		return result.Ok(ret)
	}

	r := l.read()
	if r.IsError(io.EOF) {
		return result.Ok(&token.Token{
			Type:     token.EndOfStream,
			Position: l.position,
			Line:     l.line,
			Column:   l.column,
		})
	}
	ch := r.Unwrap()

	switch {
	case ch == '{':
		return result.Ok(l.token(ch, token.OpenBrace))
	case ch == '}':
		return result.Ok(l.token(ch, token.CloseBrace))
	case ch == '(':
		return result.Ok(l.token(ch, token.OpenParen))
	case ch == ')':
		return result.Ok(l.token(ch, token.CloseParen))
	case ch == ',':
		return result.Ok(l.token(ch, token.Comma))
	case ch == ':':
		return result.Ok(l.token(ch, token.Colon))
	case ch == '@':
		return result.Ok(l.token(ch, token.At))
	case isCharacter(ch):
		builder := &strings.Builder{}
		builder.WriteRune(ch)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
			Type:     token.String,
		}
		l.consume(ch)
		for {
			r := l.read()
			if r.IsError(io.EOF) {
				tok.Capture = builder.String()
				return result.Ok(tok)
			}
			ch = r.Unwrap()
			if isCharacter(ch) {
				builder.WriteRune(ch)
				l.consume(ch)
			} else {
				l.unread().Unwrap()
				tok.Capture = builder.String()
				return result.Ok(tok)
			}
		}
	case isWhitespace(ch):
		builder := &strings.Builder{}
		builder.WriteRune(ch)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
			Type:     token.Whitespace,
		}
		l.consume(ch)
		for {
			r := l.read()
			if r.IsError(io.EOF) {
				tok.Capture = builder.String()
				return result.Ok(tok)
			}
			ch = r.Unwrap()
			if isWhitespace(ch) {
				builder.WriteRune(ch)
				l.consume(ch)
			} else {
				l.unread()
				tok.Capture = builder.String()
				return result.Ok(tok)
			}
		}
	case ch == '/':
		builder := &strings.Builder{}
		builder.WriteRune(ch)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
		}
		l.consume(ch)

		ch := l.read().Unwrap()
		builder.WriteRune(ch)
		l.consume(ch)

		switch ch {
		case '/':
			tok.Type = token.LineComment
			for {
				r := l.read()
				if r.IsError(io.EOF) {
					break
				}
				if r.IsError() {
					break
				}
				ch = r.Unwrap()
				if ch == '\n' {
					_ = l.unread().Unwrap()
					break
				}
				l.consume(ch)
				builder.WriteRune(ch)
			}
			tok.Capture = builder.String()
			return result.Ok(tok)

		case '*':
			tok.Type = token.BlockComment
			for {
				ch = l.read().Unwrap()
				builder.WriteRune(ch)
				l.consume(ch)
				if ch != '*' {
					continue
				}
				ch = l.read().Unwrap()
				builder.WriteRune(ch)
				l.consume(ch)
				if ch != '/' {
					continue
				}
				break
			}
			tok.Capture = builder.String()
			return result.Ok(tok)

		default:
			return result.Errorf[*token.Token]("invalid comment found %c expected '/' or '*'", r)
		}
	}

	return result.Errorf[*token.Token]("unrecognized token %c", r)
}

func (l *lexer) token(ch rune, ty token.TokenType) *token.Token {
	tok := &token.Token{
		Position: l.position,
		Column:   l.column,
		Line:     l.line,
		Capture:  string(ch),
		Type:     ty,
	}
	l.consume(ch)
	return tok
}

func (l *lexer) Peek() (*token.Token, error) {
	return l.peek().Deconstruct()
}

func (l *lexer) peek() (res types.Result[*token.Token]) {
	if l.peekTok != nil {
		result.Ok(l.peek)
	}
	return l.next()
}

func (l *lexer) unread() types.Result[any] {
	return result.New[any](nil, l.reader.UnreadRune())
}

func (l *lexer) read() types.Result[rune] {
	r, _, err := l.reader.ReadRune()
	return result.New(r, err)
}

func (l *lexer) consume(r rune) {
	l.position++
	l.column++
	if r == '\n' {
		l.line++
		l.column = 0
	}
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\f' || r == '\r' || r == '\n'
}

func isCharacter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '.' || r == '$'
}
