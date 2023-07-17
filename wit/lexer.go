package wit

import (
	"bufio"
	"errors"
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
		position: -1,
		line:     -1,
		column:   -1,
	}
}

func (l *lexer) Next() (*token.Token, error) {
	if l.peek != nil {
		ret := l.peek
		l.peek = nil
		return ret, nil
	}

	state := None
	capture := strings.Builder{}
	tok := &token.Token{
		Position: l.position,
		Column:   l.column,
		Line:     l.line,
	}

	for {
		r, _, err := l.reader.ReadRune()
		if err != nil && errors.Is(err, io.EOF) {
			// at end of file, cleanup
			switch state {
			case None:
				tok.Type = token.EndOfStream
				return tok, nil
			case String:
				tok.Type = token.String
			case WhiteSpace:
				tok.Type = token.Whitespace
				tok.Capture = capture.String()
				return tok, nil
			case LineComment:
				tok.Type = token.LineComment
				tok.Capture = capture.String()
				return tok, nil
			}
		}
		if err != nil {
			return nil, err
		}

		switch state {
		case None:
			switch {
			case isWhitespace(r):
				capture.WriteRune(r)
				state = WhiteSpace
			case r == '/':
				capture.WriteRune(r)
				state = BeginComment
			case r == '\n':
				capture.WriteRune(r)
				state = WhiteSpace
			}
		case WhiteSpace:
			switch {
			case isWhitespace(r) || r == '\n':
				capture.WriteRune(r)
			default:
				l.unread(r)
				tok.Capture = capture.String()
				tok.Type = token.Whitespace
			}
		case BeginComment:
			switch {
			case r == '/':
				capture.WriteRune(r)
				state = LineComment
			case r == '*':
				capture.WriteRune(r)
				state = BlockComment
			}
		case BlockComment:
			switch {
			case r == '*':
				capture.WriteRune(r)
				state = BlockCommentStar
			default:
				capture.WriteRune(r)
			}
		case BlockCommentStar:
			switch {
			case r == '/':
				capture.WriteRune(r)
				tok.Capture = capture.String()
				tok.Type = token.BlockComment
				return tok, nil
			default:
				capture.WriteRune(r)
				state = BlockComment
				capture.WriteRune(r)
			}
		case LineComment:
			switch {
			case r == '\n':
				capture.WriteRune(r)
				tok.Capture = capture.String()
				tok.Type = token.LineComment
				return tok, nil
			default:
				capture.WriteRune(r)
			}
		case String:
		}
		l.consume(r)
	}
}

func (l *lexer) Peek() (*token.Token, error) {
	if l.peek != nil {
		return l.peek, nil
	}
	var err error
	l.peek, err = l.Next()
	return l.peek, err
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

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\f' || r == '\r'
}

func isCharacter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && 'Z' <= r || '0' <= r && r <= '9' || r == '.' || r == '$'
}
