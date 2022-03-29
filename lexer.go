package wasm

import (
	"bufio"
	"io"
	"strings"
)

type lexer struct {
	reader   *bufio.Reader
	peek     *Token
	position int
	line     int
	column   int
}

type Lexer interface {
	Next() (*Token, error)
	Peek() (*Token, error)
}

func NewLexer[T string | *bufio.Reader](input T) Lexer {
	var reader *bufio.Reader
	switch v := any(input).(type) {
	case string:
		reader = bufio.NewReader(strings.NewReader(v))
	case *bufio.Reader:
		reader = v
	}
	return &lexer{
		reader: reader,
	}
}

type TokenType string

const (
	None        TokenType = "nil"
	OpenParen   TokenType = "("
	CloseParen  TokenType = ")"
	String      TokenType = "\\w+"
	Whitespace  TokenType = "\\s+"
	EndOfStream TokenType = "EOF"
)

type Token struct {
	Type     TokenType
	Position int
	Column   int
	Line     int
	Capture  string
}

func (l *lexer) Peek() (*Token, error) {
	if l.peek != nil {
		return l.peek, nil
	}
	var err error
	l.peek, err = l.Next()
	return l.peek, err
}

func (l *lexer) Next() (*Token, error) {

	if l.peek != nil {
		ret := l.peek
		l.peek = nil
		return ret, nil
	}

	state := None
	capture := strings.Builder{}
	token := &Token{
		Position: l.position,
		Column:   l.column,
		Line:     l.line,
	}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				switch state {
				case None:
					token.Type = EndOfStream
					return token, nil

				case String:
					token.Type = String
					token.Capture = capture.String()
					return token, nil

				case Whitespace:
					token.Type = Whitespace
					token.Capture = capture.String()
					return token, nil
				}
			}
			return nil, err
		}

		switch {

		case l.isWhitespace(r):
			switch state {
			case None:
				capture.WriteRune(r)
				state = Whitespace
			case Whitespace:
				capture.WriteRune(r)
			case String:
				token.Type = String
				token.Capture = capture.String()
				return token, l.unread(r)
			}

		case l.isCharacter(r):
			switch state {
			case None:
				capture.WriteRune(r)
				state = String
			case String:
				capture.WriteRune(r)
			case Whitespace:
				token.Type = Whitespace
				token.Capture = capture.String()
				return token, l.unread(r)
			}

		case r == '(':
			switch state {
			case None:
				token.Type = OpenParen
				token.Capture = string(OpenParen)
				return token, l.consume(r)

			case String:
				token.Type = String
				token.Capture = capture.String()
				return token, l.unread(r)

			case Whitespace:
				token.Type = Whitespace
				token.Capture = capture.String()
				return token, l.unread(r)
			}

		case r == ')':
			switch state {
			case None:
				token.Type = CloseParen
				token.Capture = ")"
				return token, l.consume(r)

			case String:
				token.Type = String
				token.Capture = capture.String()
				return token, l.unread(r)

			case Whitespace:
				token.Type = Whitespace
				token.Capture = capture.String()
				return token, l.unread(r)
			}
		}
		l.consume(r)
	}
}

func (l *lexer) unread(r rune) error {
	return l.reader.UnreadRune()
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

func (l *lexer) isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\f'
}

func (l *lexer) isCharacter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && 'Z' <= r || '0' <= r && r <= '9' || r == '.' || r == '$'
}
