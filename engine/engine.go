package engine

import (
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/module"
	"github.com/patrickhuber/go-wasm/store"
)

type Engine struct {
	Store *store.Store
}

func New(s *store.Store) *Engine {
	return &Engine{}
}

func (e *Engine) Start(m module.Module) *instance.Module {

	return &instance.Module{}
}
