package lex

import "github.com/patrickhuber/go-wasm/wat/token"

const RuneSliceRuleType RuleType = 0

type RuneSliceRule struct {
	Runes     []rune
	TokenType token.Type
}

func (r *RuneSliceRule) Type() token.Type {
	return r.TokenType
}

func (r *RuneSliceRule) Check(ch rune) bool {
	if len(r.Runes) == 0 {
		return false
	}
	return r.Runes[0] == ch
}

func (r *RuneSliceRule) RuleType() RuleType {
	return RuneSliceRuleType
}

type RuneSliceLexeme struct {
	position int
	rule     *RuneSliceRule
}

func (l *RuneSliceLexeme) Rule() Rule {
	return l.rule
}

func (r *RuneSliceLexeme) Eat(ch rune) bool {
	if r.position >= len(r.rule.Runes) {
		return false
	}
	if r.rule.Runes[r.position] != ch {
		return false
	}
	r.position++
	return true
}

func (r *RuneSliceLexeme) Accepted() bool {
	return r.position >= len(r.rule.Runes)
}

type RuneSliceFactory struct {
}

func (f *RuneSliceFactory) Lexeme(r Rule) Lexeme {
	rule := r.(*RuneSliceRule)
	return &RuneSliceLexeme{
		position: 0,
		rule:     rule,
	}
}
