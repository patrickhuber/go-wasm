package wat

import (
	"bufio"
	"fmt"
	"io"

	"github.com/patrickhuber/go-collections/generic/stack"
)

type Reader interface {
	Next() (bool, error)
	Current() Node
}

type reader struct {
	lexer   Lexer
	current Node
	stack   stack.Stack[Node]
}

type readError struct {
	token   *Token
	message string
	args    []any
}

func (err *readError) Error() string {
	msg := fmt.Sprintf(err.message, err.args...)
	return fmt.Sprintf("parse error line: %d, column: %d, position: %d. %s", err.token.Line+1, err.token.Column+1, err.token.Position, msg)
}

type Node interface {
	Type() NodeType
}

type NodeType string

const (
	ModuleNodeType   NodeType = "module"
	ResultNodeType   NodeType = "result"
	MemoryNodeType   NodeType = "memory"
	FuncNodeType     NodeType = "func"
	ArgumentNodeType NodeType = "arg"
	ParamNodeType    NodeType = "param"
	LocalNodeType    NodeType = "local"
	ExportNodeType   NodeType = "export"
	EndNodeType      NodeType = "end"
)

type node struct {
	nodeType NodeType
	value    string
}

func (n *node) Type() NodeType {
	return n.nodeType
}

func NewNode(nodeType NodeType) Node {
	return &node{nodeType: nodeType}
}

func NewReader(r io.Reader) Reader {
	br := bufio.NewReader(r)
	return &reader{
		lexer: NewLexer(br),
		stack: stack.New[Node](),
	}
}

func (r *reader) Current() Node {
	return r.current
}

func (r *reader) Next() (bool, error) {
	t, err := r.peek()
	if err != nil {
		return false, err
	}
	switch t.Type {
	case OpenParen:
		n, err := r.node()
		if err != nil {
			return false, err
		}
		r.current = n
		r.stack.Push(n)
		return true, nil

	case CloseParen:
		n, err := r.end()
		if err != nil {
			return false, err
		}
		if n.Type() != EndNodeType {
			return false, r.error(t, "expected end node")
		}
		if r.stack.Length() == 0 {
			return false, r.error(t, "mismatched end node")
		}
		r.stack.Pop()
		r.current = n
		return true, nil

	case String:
		n, err := r.arg()
		if err != nil {
			return false, err
		}
		r.current = n
		return true, nil
	case EndOfStream:
		return false, nil
	}

	return false, r.error(t, "unexpected token type: '%s' value: '%s'", t.Type, t.Capture)
}

func (r *reader) node() (Node, error) {
	err := r.expectToken(OpenParen)
	if err != nil {
		return nil, err
	}
	t, err := r.string()
	if err != nil {
		return nil, err
	}
	return &node{nodeType: NodeType(t.Capture)}, nil
}

func (r *reader) arg() (Node, error) {
	t, err := r.string()
	if err != nil {
		return nil, err
	}
	return &node{nodeType: NodeType(ArgumentNodeType), value: t.Capture}, nil
}

func (r *reader) end() (Node, error) {
	err := r.expectToken(CloseParen)
	if err != nil {
		return nil, err
	}
	return &node{nodeType: EndNodeType}, nil
}

func (r *reader) string() (*Token, error) {
	token, err := r.nextToken()
	if err != nil {
		return nil, err
	}
	if token.Type != String {
		return nil, r.error(token, "expected '%s' found '%s' ", String, token.Type)
	}
	return token, nil
}

func (r *reader) expectToken(t TokenType) error {
	token, err := r.nextToken()
	if err != nil {
		return err
	}
	if token.Type != t {
		return r.error(token, "expected '%s' found '%s'", t, token.Type)
	}
	return nil
}

func (r *reader) error(t *Token, msg string, args ...any) error {
	return &readError{
		token:   t,
		message: msg,
		args:    args,
	}
}

func (r *reader) peek() (*Token, error) {

	for {
		t, err := r.lexer.Peek()
		if err != nil {
			return nil, err
		}

		// return peeked token
		if t.Type != Whitespace {
			return t, nil
		}

		// consume the whitespace
		r.lexer.Next()
	}
}

func (r *reader) nextToken() (*Token, error) {
	_, err := r.peek()
	if err != nil {
		return nil, err
	}
	return r.lexer.Next()
}
