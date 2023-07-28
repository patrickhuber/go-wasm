package lex

import (
	"github.com/patrickhuber/go-wasm/wat/token"
)

type Rule interface {
	Check(r rune) bool
	Type() token.Type
}

type Lexeme interface {
	Eat(ch rune) bool
	Rule() Rule
}

type RuneSliceLexeme struct {
	runes    []rune
	position int
	rule     Rule
}

func (r *RuneSliceLexeme) Eat(ch rune) bool {
	if r.position >= len(r.runes) {
		return false
	}
	if r.runes[r.position] != ch {
		return false
	}
	r.position++
	return true
}

func (r *RuneSliceLexeme) Rule() Rule {
	return r.rule
}

type RuneLexeme struct {
	ch   rune
	done bool
	rule Rule
}

func (r *RuneLexeme) Eat(ch rune) bool {
	if r.done {
		return false
	}
	if r.ch == ch {
		r.done = true
		return true
	}
	return false
}

func (r *RuneLexeme) Rule() Rule {
	return r.rule
}
