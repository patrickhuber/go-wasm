package ast

import "github.com/patrickhuber/go-types"

type Wat interface {
	wat()
}

type Component struct {
	Wat
}

type Module struct {
	Wat
	Functions []Function
	Memory    []Memory
	Types     []Type
	Tables    []Table
	Globals   []Global
}

type Section interface {
	section()
}

type Function struct {
	ID           types.Option[string]
	Locals       []Local
	Exports      []InlineExport
	Parameters   []Parameter
	Results      []Result
	Instructions []Instruction
}

func (f *Function) section() {}

type Memory struct {
	ID     types.Option[string]
	Limits Limits
}

type Type struct {
	ID       types.Option[string]
	FuncType FuncType
}

type FuncType struct {
	Parameters []Parameter
	Results    []Result
}

type Table struct {
	ID        types.Option[string]
	TableType TableType
	Elements  []Element
}

type TableType struct {
	Limits  Limits
	RefType RefType
}

type Global struct {
	ID           types.Option[string]
	Type         GlobalType
	Instructions []Instruction
}

type GlobalType struct {
	Mutable bool
	Type    ValType
}

type Element struct {
	ID string
}

type Limits struct {
	Min uint32
	Max types.Option[uint32]
}

type Parameter struct {
	ID    types.Option[string]
	Types []ValType
}

type Local struct {
	Type ValType
}

type Result struct {
	Types []ValType
}

type RefType interface {
	refType()
}

type FuncRef struct{}

func (FuncRef) refType() {}

type ExternRef struct{}

func (ExternRef) refType() {}

type ValType interface {
	valType()
}

type I32 struct{}

func (I32) valType() {}

type I64 struct{}

func (I64) valType() {}

type F32 struct{}

func (F32) valType() {}

type F64 struct{}

func (F64) valType() {}

type Instruction interface {
	inst()
}

type I32Const struct {
	Value int32
}

func (I32Const) inst() {}

type I64Const struct {
	Value int64
}

func (I64Const) inst() {}

type F32Const struct {
	Value float32
}

func (F32Const) inst() {}

type F64Const struct {
	Value float64
}

func (F64Const) inst() {}

type I32Eqz struct {
	Value int32
}

func (I32Eqz) inst() {}

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

type LocalSet struct {
	Index Index
}

func (LocalSet) inst() {}

type LocalTee struct {
	Index Index
}

func (LocalTee) inst() {}

type GlobalGet struct {
	Index Index
}

func (GlobalGet) inst() {}

type GlobalSet struct {
	Index Index
}

func (GlobalSet) inst() {}

type MemoryGrow struct{}

func (MemoryGrow) inst() {}

type I32Load struct{}

func (I32Load) inst() {}

type I32Store struct{}

func (I32Store) inst() {}

type Drop struct{}

func (Drop) inst() {}

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

type Block struct {
	Name         types.Option[string]
	BlockType    BlockType
	Instructions []Instruction
}

func (Block) inst() {}

type BlockType struct {
	Results []Result
}

type Loop struct {
	Name         types.Option[string]
	BlockType    BlockType
	Instructions []Instruction
}

func (Loop) inst() {}

type If struct {
	Name      types.Option[string]
	Clause    []Instruction
	BlockType BlockType
	Then      Then
	Else      types.Option[Else]
}

func (If) inst() {}

type Then struct {
	Instructions []Instruction
}

type Else struct {
	Instructions []Instruction
}

type Br struct {
	Index Index
}

func (Br) inst() {}

type BrIf struct {
	Index Index
}

func (BrIf) inst() {}

type BrTable struct {
	Indicies []Index
}

func (BrTable) inst() {}

type Return struct{}

func (Return) inst() {}

type Select struct{}

func (Select) inst() {}

type Call struct {
	Index Index
}

func (Call) inst() {}

type CallIndirect struct {
	Type TypeUse
}

func (CallIndirect) inst() {}

type TypeUse struct {
	Index string
}
