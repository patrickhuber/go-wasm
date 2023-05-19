package runtime

import (
	"github.com/patrickhuber/go-wasm/store"
)

type Engine interface {
	Instantiate() error
	Invoke() error
}
type engine struct {
	stack Stack
	store store.Store
}

func NewEngine(stack Stack, store store.Store) Engine {
	return &engine{
		stack: stack,
		store: store,
	}
}

func (e *engine) Instantiate() error {
	return nil
}

func (e *engine) Invoke() error {
	return nil
}

func (e *engine) I32Const(value uint32) {
	e.stack.PushUint32(value)
}

func (e *engine) I32Add() {
	i1 := e.stack.PopUint32()
	i2 := e.stack.PopUint32()
	e.stack.PushUint32(i1 + i2)
}
