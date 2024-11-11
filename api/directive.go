package api

type Directive interface {
	directive()
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
	Body   Expression
}

type Mem struct {
}

type Import struct{}
type Export struct{}
type Start struct{}
type Data struct{}
type Elem struct{}
