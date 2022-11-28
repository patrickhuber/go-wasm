package runtime

type Store struct {
	Functions []FunctionInstance
	Tables    []TableInstance
	Memories  []MemoryInstance
	Globals   []GlobalInstance
	Elements  []ElementInstance
	Datas     []DataInstance
}

type FunctionInstance struct{}
type TableInstance struct{}

type MemoryInstance struct {
	Max          int
	MaxSpecified bool
	Data         []byte
}

type GlobalInstance struct {
	Value   Value
	Mutable bool
}
type ElementInstance struct{}
type DataInstance struct{}
