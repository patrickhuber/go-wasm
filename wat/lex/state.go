package lex

import "github.com/patrickhuber/go-wasm/wat/token"

const StateRuleType RuleType = 3

type StateRule struct {
	TokenType     token.Type
	StartMatch    func(b byte) bool
	ContinueMatch func(state string, b byte) (string, bool)
}

func (s *StateRule) Type() token.Type {
	return s.TokenType
}

func (s *StateRule) Check(ch byte) bool {
	return s.StartMatch(ch)
}

func (*StateRule) RuleType() RuleType {
	return StateRuleType
}

type StateLexeme struct {
	rule    *StateRule
	current string
}

func (l *StateLexeme) Eat(ch byte) bool {
	if l.current == "" {
		return l.rule.StartMatch(ch)
	}
	current, ok := l.rule.ContinueMatch(l.current, ch)
	if ok {
		l.current = current
	}
	return ok
}
