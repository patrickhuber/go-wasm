package runtime

import (
	"encoding/binary"

	"github.com/patrickhuber/go-wasm/wasm"
)

type engine struct {
	stack Stack
	store Store
}

func (e *engine) I32Const(value uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	item := StackItem{
		Value: &Value{
			Number: &Number{
				OpCode: wasm.I32Const,
				Value:  b,
			},
		},
	}
	e.stack.Push(item)
}

func (e *engine) I32Add() {
	i1 := e.stack.Pop()
	i2 := e.stack.Pop()
	value1 := binary.BigEndian.Uint32(i1.Value.Number.Value)
	value2 := binary.BigEndian.Uint32(i2.Value.Number.Value)
	binary.BigEndian.PutUint32(i1.Value.Number.Value, value1+value2)

	e.stack.Push(i1)
}

func (e *engine) popI32() uint32 {
	item := e.stack.Pop()
	number := item.Value.Number
	return binary.BigEndian.Uint32(number.Value)
}

func (e *engine) pushI32(value uint32) {
	// this memory needs to come from the store?
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	item := StackItem{
		Value: &Value{
			Number: &Number{
				OpCode: wasm.I32Const,
				Value:  b,
			},
		},
	}
	e.stack.Push(item)
}
