package runtime

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/address"
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/values"
)

// Machine models execution behavior in terms of an abstract machine that models the program state.
type Machine struct {
	store *Store
	stack *Stack
}

func NewMachine() *Machine {
	return &Machine{
		store: &Store{},
		stack: &Stack{},
	}
}

// Execute executes a wasm program
func (m *Machine) Execute(module *api.Module, externals []values.Value) ([]values.Value, error) {
	err := m.validate(module, externals)
	if err != nil {
		return nil, err
	}
	err = m.allocate(module)
	if err != nil {
		return nil, err
	}
	return []values.Value{}, nil
}

func (m *Machine) validate(module *api.Module, externals []values.Value) error {
	if len(externals) != len(module.Imports) {
		return fmt.Errorf("module imports should match external values")
	}
	return nil
}

func (m *Machine) allocate(module *api.Module) error {
	moduleInstance := &instance.Module{}
	for _, fn := range module.Funcs {
		funcAddr, err := m.allocFunc(moduleInstance, fn)
		if err != nil {
			return err
		}
		moduleInstance.FunctionAddresses = append(moduleInstance.FunctionAddresses, funcAddr)
	}
	return nil
}

// https://webassembly.github.io/spec/core/exec/modules.html#functions
func (m *Machine) allocFunc(moduleInst *instance.Module, fn *api.Func) (address.Function, error) {
	// first free address
	funcAddress := len(m.store.Funcs)
	funcType := moduleInst.Types[fn.Type]
	funcInstance := instance.ModuleFunction{
		Type:   funcType,
		Code:   fn,
		Module: moduleInst,
	}
	m.store.Funcs = append(m.store.Funcs, funcInstance)
	return address.Function{
		Address: uint32(funcAddress),
	}, nil
}
