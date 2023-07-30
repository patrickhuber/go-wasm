package lex

import "github.com/patrickhuber/go-wasm/wat/token"

const StringRuleType RuleType = 0

type StringRule struct {
	String    string
	TokenType token.Type
}

func (r *StringRule) Type() token.Type {
	return r.TokenType
}

func (r *StringRule) Check(ch byte) bool {
	if len(r.String) == 0 {
		return false
	}
	return r.String[0] == ch
}

func (r *StringRule) RuleType() RuleType {
	return StringRuleType
}

type StringLexeme struct {
	position int
	rule     *StringRule
}

func (l *StringLexeme) Rule() Rule {
	return l.rule
}

func (r *StringLexeme) Eat(ch byte) bool {
	if r.position >= len(r.rule.String) {
		return false
	}
	if r.rule.String[r.position] != ch {
		return false
	}
	r.position++
	return true
}

func (r *StringLexeme) Accepted() bool {
	return r.position >= len(r.rule.String)
}

type StringFactory struct {
}

func (f *StringFactory) Lexeme(r Rule) Lexeme {
	rule := r.(*StringRule)
	return &StringLexeme{
		position: 0,
		rule:     rule,
	}
}
