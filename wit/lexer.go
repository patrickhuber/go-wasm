package wit

import "bufio"

type TokenType string

const (
	None        TokenType = "nil"
	OpenParen   TokenType = "("
	CloseParen  TokenType = ")"
	String      TokenType = "\\w+"
	Whitespace  TokenType = "\\s+"
	EndOfStream TokenType = "EOF"
)

type Lexer interface {
	Next() (*Token, error)
	Peek() (*Token, error)
}

type lexer struct {
	reader   *bufio.Reader
	peek     *Token
	position int
	line     int
	column   int
}

type Token struct {
	Type     TokenType
	Position int
	Column   int
	Line     int
	Capture  string
}
