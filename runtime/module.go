package runtime

type ModuleInstance struct {
	Types             []FunctionType
	FunctionAddresses []FunctionAddress
	TableAddresses    []TableAddress
	MemoryAddresses   []MemoryAddress
	GlobalAddresses   []GlobalAddress
	ElementAddresses  []ElementAddress
	DataAddresses     []DataAddress
	Exports           []ExportInstance
}

type FunctionType struct{}
type FunctionAddress struct{}
type TableAddress struct{}
type MemoryAddress struct{}
type GlobalAddress struct{}
type ElementAddress struct{}
type DataAddress struct{}
type ExportInstance struct{}
