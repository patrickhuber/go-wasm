package api

type Directive interface {
	directive()
}

type Document struct {
	Preamble  Preamble
	Directive Directive
}

const (
	ModuleVersion    uint16 = 1
	ComponentVersion uint16 = 2
)

type Preamble struct {
	Version uint16
	Layer   uint16
}

type Module struct {
	Types   []*FuncType
	Funcs   []*Func
	Tables  []Table
	Mems    []Mem
	Globals []Global
	Elems   []Elem
	Datas   []Data
	Start   Start
	Imports []Import
	Exports []Export
}

func (*Module) directive() {}

type Component struct{}

func (*Component) directive() {}

type FuncType struct {
	Parameters ResultType
	Returns    ResultType
}

type Func struct {
	Type   TypeIndex
	Locals []ValType
	Body   *Expression
}

type Mem struct {
}

type Import struct{}
type Export struct {
	Name        string
	Description ExportDescription
}
type ExportDescription interface {
	exportDescription()
}
type FuncExportDescription struct {
	FuncIdx FuncIndex
}

func (*FuncExportDescription) exportDescription() {}

type Start struct{}
type Data struct{}
type Elem struct{}
