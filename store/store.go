package store

import (
	"github.com/patrickhuber/go-wasm/engine"
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/wasm"
)

func New(engine *engine.Engine) *Store {
	return &Store{
		Engine: engine,
	}
}

type Store struct {
	Engine    *engine.Engine
	Directive instance.Directive
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
