package ast

import wat "github.com/patrickhuber/go-wasm/wat/ast"

type Directive interface {
	directive()
}

type AssertReturn struct{}

func (AssertReturn) directive() {}

type Wat struct {
	Wat wat.Ast
}

func (Wat) directive() {}

type AssertInvalid struct{}

func (AssertInvalid) directive() {}

type AssertMalformed struct{}

func (AssertMalformed) directive() {}

type Invoke struct{}

func (Invoke) directive() {}

type AssertTrap struct{}

func (AssertTrap) directive() {}
