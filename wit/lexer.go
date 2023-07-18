package wit

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/patrickhuber/go-wasm/wit/token"
)

type Lexer interface {
	Next() (*token.Token, error)
	Peek() (*token.Token, error)
}

type lexer struct {
	reader   *bufio.Reader
	peek     *token.Token
	position int
	line     int
	column   int
}

func NewLexer(reader io.Reader) Lexer {
	return &lexer{
		reader:   bufio.NewReader(reader),
		peek:     nil,
		position: 0,
		line:     0,
		column:   0,
	}
}

func (l *lexer) Next() (*token.Token, error) {
	if l.peek != nil {
		ret := l.peek
		l.peek = nil
		return ret, nil
	}

	r, _, err := l.reader.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return &token.Token{
				Type:     token.EndOfStream,
				Position: l.position,
				Line:     l.line,
				Column:   l.column,
			}, nil
		}
		return nil, err
	}

	switch {
	case r == '{':
		return l.token(r, token.OpenBrace), nil
	case r == '}':
		return l.token(r, token.CloseBrace), nil
	case r == '(':
		return l.token(r, token.OpenParen), nil
	case r == ')':
		return l.token(r, token.CloseParen), nil
	case r == ',':
		return l.token(r, token.Comma), nil
	case r == ':':
		return l.token(r, token.Colon), nil
	case r == '@':
		return l.token(r, token.At), nil
	case isCharacter(r):
		builder := &strings.Builder{}
		builder.WriteRune(r)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
			Type:     token.String,
		}
		l.consume(r)
		for {
			r, err := l.read()
			if err != nil && errors.Is(err, io.EOF) {
				tok.Capture = builder.String()
				return tok, nil
			}
			if err != nil {
				return nil, err
			}
			if isCharacter(r) {
				builder.WriteRune(r)
				l.consume(r)
			} else {
				l.unread()
				tok.Capture = builder.String()
				return tok, nil
			}
		}
	case isWhitespace(r):
		builder := &strings.Builder{}
		builder.WriteRune(r)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
			Type:     token.Whitespace,
		}
		l.consume(r)
		for {
			r, err := l.read()
			if err != nil && errors.Is(err, io.EOF) {
				tok.Capture = builder.String()
				return tok, nil
			}
			if err != nil {
				return nil, err
			}
			if isWhitespace(r) {
				builder.WriteRune(r)
				l.consume(r)
			} else {
				l.unread()
				tok.Capture = builder.String()
				return tok, nil
			}
		}
	case r == '/':
		builder := &strings.Builder{}
		builder.WriteRune(r)
		tok := &token.Token{
			Position: l.position,
			Column:   l.column,
			Line:     l.line,
		}
		l.consume(r)

		r, err := l.read()
		if err != nil {
			return nil, err
		}
		builder.WriteRune(r)
		l.consume(r)

		switch r {
		case '/':
			tok.Type = token.LineComment
			for {
				r, err := l.read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break
					}
					return nil, err
				}
				if r == '\n' {
					l.unread()
					break
				}
				l.consume(r)
				builder.WriteRune(r)
			}
			tok.Capture = builder.String()
			return tok, nil

		case '*':
			tok.Type = token.BlockComment
			for {
				r, err := l.read()
				if err != nil {
					return nil, err
				}
				builder.WriteRune(r)
				l.consume(r)
				if r != '*' {
					continue
				}
				r, err = l.read()
				if err != nil {
					return nil, err
				}
				builder.WriteRune(r)
				l.consume(r)
				if r != '/' {
					continue
				}
				break
			}
			tok.Capture = builder.String()
			return tok, nil

		default:
			return nil, fmt.Errorf("invalid comment found %c expected '/' or '*'", r)
		}
	}

	return nil, fmt.Errorf("unrecognized token %c", r)
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
	if l.peek != nil {
		return l.peek, nil
	}
	var err error
	l.peek, err = l.Next()
	return l.peek, err
}

func (l *lexer) unread() error {
	return l.reader.UnreadRune()
}

func (l *lexer) read() (rune, error) {
	r, _, err := l.reader.ReadRune()
	return r, err
}

func (l *lexer) consume(r rune) error {
	l.position++
	l.column++
	if r == '\n' {
		l.line++
		l.column = 0
	}
	return nil
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\f' || r == '\r' || r == '\n'
}

func isCharacter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '.' || r == '$'
}
