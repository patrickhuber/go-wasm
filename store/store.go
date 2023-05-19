package store

import "github.com/patrickhuber/go-wasm/wasm"

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
	Value   wasm.Value
	Mutable bool
}
type ElementInstance struct{}
type DataInstance struct{}
