package ast

import (
	"github.com/patrickhuber/go-types"
	wat "github.com/patrickhuber/go-wasm/wat/ast"
)

type Directive interface {
	directive()
}

type AssertReturn struct {
	Action  Action
	Results []Result
}

func (AssertReturn) directive() {}

type Wat struct {
	Wat wat.Ast
}

func (Wat) directive() {}

type AssertInvalid struct {
	Module  *wat.Module
	Failure string
}

func (AssertInvalid) directive() {}

type AssertMalformed struct{}

func (AssertMalformed) directive() {}

type AssertTrap struct {
	Action  Action
	Failure string
}

func (AssertTrap) directive() {}

type Action interface {
	action()
}

type Invoke struct {
	Name   types.Option[string]
	String string
	Const  []Const
}

func (Invoke) action() {}

type Get struct{}

func (Get) action() {}

type Const interface {
	const_()
}

type Result interface {
	result()
}

type I32Const struct {
	Value int32
}

func (I32Const) const_() {}
func (I32Const) result() {}

type I64Const struct {
	Value int64
}

func (I64Const) const_() {}
func (I64Const) result() {}

type F32Const struct {
	Value float32
}

func (F32Const) const_() {}
func (F32Const) result() {}

type F64Const struct {
	Value float64
}

func (F64Const) const_() {}
func (F64Const) result() {}
