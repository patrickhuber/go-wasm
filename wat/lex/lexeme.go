package lex

import (
	"github.com/patrickhuber/go-wasm/wat/token"
)

type RuleType int

type Rule interface {
	Check(r rune) bool
	Type() token.Type
	RuleType() RuleType
}

type Lexeme interface {
	Eat(ch rune) bool
	Rule() Rule
	Accepted() bool
}
