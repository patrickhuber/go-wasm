package lex

import (
	"github.com/patrickhuber/go-wasm/wat/token"
)

type RuleType int

type Rule interface {
	Check(r byte) bool
	Type() token.Type
	RuleType() RuleType
}

type Lexeme interface {
	Eat(ch byte) bool
	Rule() Rule
	Accepted() bool
}
