package stack

import (
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/instruction"
	"github.com/patrickhuber/go-wasm/values"
)

type Item interface {
	item()
}

type Label struct {
	Instructions []instruction.Instruction
}

func (*Label) item() {}

type Value struct{}

func (*Value) item() {}

type Frame struct {
	Locals values.Value
	Module instance.Module
}

func (*Frame) item() {}
