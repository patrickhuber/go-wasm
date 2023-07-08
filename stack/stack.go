package stack

import (
	"github.com/patrickhuber/go-wasm/instances"
	"github.com/patrickhuber/go-wasm/instructions"
	"github.com/patrickhuber/go-wasm/values"
)

type Item interface {
	item()
}

type Label struct {
	Instructions []instructions.Instruction
}

func (*Label) item() {}

type Value struct{}

func (*Value) item() {}

type Frame struct {
	Locals values.Value
	Module instances.Module
}

func (*Frame) item() {}
