package ast

import (
	"github.com/patrickhuber/go-types"
	wat "github.com/patrickhuber/go-wasm/wat/ast"
)

type Directive interface {
	directive()
}

type Wast struct {
	Directives []Directive
}

type QuoteWat interface {
	quoteWat()
}

type WatDirective struct {
	Directive
	Wat QuoteWat
}

type Wat struct {
	QuoteWat
	Wat wat.Directive
}

type QuoteModule struct {
	QuoteWat
	Quote string
}

type QuoteComponent struct {
	QuoteWat
	Quote string
}

type AssertInvalid struct {
	Directive
	Module  QuoteWat
	Failure string
}

type AssertMalformed struct {
	Directive
	Module  QuoteWat
	Failure string
}

type AssertTrap struct {
	Directive
	Action  Action
	Failure string
}

type AssertReturn struct {
	Directive
	Action  Action
	Results []Result
}

type Action interface {
	action()
}

type Invoke struct {
	Action
	Directive
	Name   types.Option[string]
	String string
	Const  []Const
}

type Get struct {
	Action
}

type Const interface {
	const_()
}

type Result interface {
	result()
}

type I32Const struct {
	Const
	Result
	Value int32
}

type I64Const struct {
	Const
	Result
	Value int64
}

type F32Const struct {
	Const
	Result
	Value float32
}

type F64Const struct {
	Const
	Result
	Value float64
}
