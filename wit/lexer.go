package wit

import (
	"bufio"

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
