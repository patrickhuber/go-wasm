package api

import (
	"github.com/patrickhuber/go-wasm/address"
	"github.com/patrickhuber/go-wasm/values"
)

type Trap struct{}

func (*Trap) instruction() {}

type Invoke struct {
	FunctionAddress address.Function
}

func (*Invoke) instruction() {}

type Label struct {
	Instructions []Instruction
	Instruction  Instruction
}

func (*Label) instruction() {}

// I don't want to create a circular reference with the
// stack.Frame struct so duplicate structure here
type Frame struct {
	Frames      []InnerFrame
	Instruction Instruction
}

func (*Frame) instruction() {}

type InnerFrame struct {
	Locals []values.Value
	Module Module
}
