package stack

import (
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/values"
)

type Item interface {
	item()
}

type Label struct {
	Instructions []api.Instruction
}

func (*Label) item() {}

type Value struct{}

func (*Value) item() {}

type Frame struct {
	Locals values.Value
	Module instance.Module
}

func (*Frame) item() {}
