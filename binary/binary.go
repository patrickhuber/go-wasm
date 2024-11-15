package binary

import (
	"github.com/patrickhuber/go-wasm/api"
)

var Magic = []byte{0x00, 0x61, 0x73, 0x6d}

type Document struct {
	Preamble  *Preamble
	Directive Directive
}

const ModuleVersion uint16 = 0x01
const ComponentVersion uint16 = 0x0a

type Preamble struct {
	Magic   []byte
	Version uint16
	Layer   uint16
}

type Directive interface {
	directive()
}

type Component struct{}

func (Component) directive() {}

type Module struct {
	Sections []Section
}

func (Module) directive() {}

type Section interface {
	section()
}

type SectionID uint8

const (
	CustomSectionID   SectionID = 0
	TypeSectionID     SectionID = 1
	FunctionSectionID SectionID = 3
	CodeSectionID     SectionID = 10
)

type TypeSection struct {
	ID    SectionID
	Size  uint32
	Types []*FunctionType
}

func (TypeSection) section() {}

type FunctionType struct {
	Parameters ResultType
	Returns    ResultType
}

type ResultType struct {
	Types []ValType
}

type FunctionSection struct {
	ID    SectionID
	Size  uint32
	Types []uint32
}

func (FunctionSection) section() {}

type CodeSection struct {
	ID    SectionID
	Size  uint32
	Codes []*Code
}

func (CodeSection) section() {}

type Code struct {
	Size       uint32
	Locals     []Local
	Expression []api.Instruction
}

type Local struct {
	ValueTypes []ValType
}

type ValType byte

const I32 ValType = 0x7f
const I64 ValType = 0x7e
const F32 ValType = 0x7d
const F64 ValType = 0x7c
