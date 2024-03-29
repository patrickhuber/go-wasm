package lex

import "github.com/patrickhuber/go-wasm/wat/token"

const DfaRuleType RuleType = 1

type Edge interface {
	Check(ch byte) bool
	Next() *Node
}

type ByteEdge struct {
	Byte byte
	Node *Node
}

func (r *ByteEdge) Check(ch byte) bool {
	return r.Byte == ch
}

func (r *ByteEdge) Next() *Node {
	return r.Node
}

type FuncEdge struct {
	Func func(ch byte) bool
	Node *Node
}

func (f *FuncEdge) Check(ch byte) bool {
	return f.Func(ch)
}

func (f *FuncEdge) Next() *Node {
	return f.Node
}

type Node struct {
	Edges []Edge
	Final bool
}

type Dfa struct {
	Start *Node
}

type DfaRule struct {
	Dfa       *Dfa
	TokenType token.Type
}

func (r *DfaRule) Type() token.Type {
	return r.TokenType
}

func (r *DfaRule) Check(ch byte) bool {
	for _, edge := range r.Dfa.Start.Edges {
		if edge.Check(ch) {
			return true
		}
	}
	return false
}

func (r *DfaRule) RuleType() RuleType {
	return DfaRuleType
}

type DfaLexeme struct {
	rule    *DfaRule
	current *Node
}

func (l *DfaLexeme) Rule() Rule {
	return l.rule
}

func (d *DfaLexeme) Eat(ch byte) bool {
	if d.current == nil {
		d.current = d.rule.Dfa.Start
	}
	// look for outbound edges
	for _, edge := range d.current.Edges {
		// dfa is deterministic so we can exit here
		if edge.Check(ch) {
			d.current = edge.Next()
			return true
		}
	}
	return false
}

func (d *DfaLexeme) Accepted() bool {
	return d.current.Final
}

type DfaFactory struct {
}

func (f *DfaFactory) Lexeme(r Rule) Lexeme {
	rule := r.(*DfaRule)
	return &DfaLexeme{
		rule: rule,
	}
}

func FromString(str string) (*Node, *Node) {
	start := &Node{}
	current := start
	for i := 0; i < len(str); i++ {
		next := &Node{}
		current.Edges = append(current.Edges, &ByteEdge{
			Byte: str[i],
			Node: next,
		})
		current = next
	}
	return start, current
}
