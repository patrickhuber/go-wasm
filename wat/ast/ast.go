package ast

import "github.com/patrickhuber/go-types"

type Ast interface {
	ast()
}

type Component struct {
}

func (*Component) ast() {}

type Module struct {
	Functions []Function
	Memory    []Function
}

func (*Module) ast() {}

type Section interface {
	section()
}

type Function struct {
	ID           types.Option[string]
	Exports      []InlineExport
	Parameters   []Parameter
	Results      []Result
	Instructions []Instruction
}

func (f *Function) section() {}

type Memory struct {
	ID     string
	Limits Limits
}

type Limits struct {
	HasMax bool
	Min    uint32
	Max    uint32
}

type Parameter struct {
	ID   types.Option[string]
	Type Type
}

type Local struct {
	Type Type
}

type Result struct {
	Type Type
}

type Type interface {
	ty()
}

type I32 struct{}

func (I32) ty() {}

type I64 struct{}

func (I64) ty() {}

type F32 struct{}

func (F32) ty() {}

type F64 struct{}

func (F64) ty() {}

type Instruction interface {
	inst()
}

type I32Add struct{}

func (I32Add) inst() {}

type I32Sub struct{}

func (I32Sub) inst() {}

type I32Mul struct{}

func (I32Mul) inst() {}

type I32DivS struct{}

func (I32DivS) inst() {}

type I32DivU struct{}

func (I32DivU) inst() {}

type Folded struct {
	Instruction Instruction
	Parameters  []Instruction
}

func (Folded) inst() {}

type LocalGet struct {
	Index Index
}

func (LocalGet) inst() {}

type Index interface {
	index()
}

type IDIndex struct {
	ID string
}

func (IDIndex) index() {}

type RawIndex struct {
	Index uint32
}

func (RawIndex) index() {}

type InlineExport struct {
	Name string
}

type InlineImport struct {
	Module string
	Field  string
}
