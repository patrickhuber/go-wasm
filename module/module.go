package module

import (
	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/instruction"
	"github.com/patrickhuber/go-wasm/types"
)

type Module struct {
	Types     []FuncType
	Functions []Function
	Memories  []Memory
	Globals   []Global
	Datas     []Data
	Start     Start
	Imports   []Import
	Exports   []Export
}

type FuncType struct {
	Parameters []types.Value
	Results    []types.Value
}

type Function struct {
	Type   indicies.Type
	Locals []types.Value
	Body   instruction.Expression
}

type Table struct {
	Type indicies.Table
}

type Memory struct {
	Type indicies.Memory
}

type Global struct {
	Type indicies.Global
	Init instruction.Expression
}

type Element struct {
	Type types.Reference
	Init []instruction.Expression
	Mode ElementMode
}

type ElementMode interface {
	elementmode()
}

type PassiveElementMode struct{}

func (*PassiveElementMode) elementmode() {}

type ActiveElementMode struct {
	Table  indicies.Table
	Offset instruction.Expression
}

func (*ActiveElementMode) elementmode() {}

type DeclaritiveElementMode struct{}

func (*DeclaritiveElementMode) elementmode() {}

type Data struct {
	Init []byte
	Mode DataMode
}

type DataMode interface {
	datamode()
}

type PassiveDataMode struct{}

func (*PassiveDataMode) elementmode() {}

type ActiveDataMode struct {
	Memory indicies.Memory
	Offset instruction.Expression
}

func (*ActiveDataMode) elementmode() {}

type Start struct {
	Func indicies.Function
}

type Export struct {
	Name        string
	Description ExportDescription
}

type ExportDescription interface {
	exportdescription()
}

type FunctionExportDescription struct {
	Func indicies.Function
}

func (*FunctionExportDescription) exportdescription() {}

type TableExportDescription struct {
	Table indicies.Table
}

func (*TableExportDescription) exportdescription() {}

type MemoryExportDescription struct {
	Memory indicies.Memory
}

func (*MemoryExportDescription) exportdescription() {}

type GlobalMemoryDescription struct {
	Global indicies.Global
}

func (*GlobalMemoryDescription) exportdescription() {}

type Import struct {
	Module      string
	Name        string
	Description ImportDescription
}

type ImportDescription interface {
	importdescription()
}

type FunctionImportDescription struct {
	Func indicies.Function
}

func (*FunctionImportDescription) importdescription() {}

type TableImportDescription struct {
	Table indicies.Table
}

func (*TableImportDescription) importdescription() {}

type MemoryImportDescription struct {
	Memory indicies.Memory
}

func (*MemoryImportDescription) importdescription() {}

type GlobalImportDescription struct {
	Global indicies.Global
}

func (*GlobalImportDescription) importdescription() {}
