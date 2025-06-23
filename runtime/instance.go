package runtime

import (
	"github.com/patrickhuber/go-wasm/address"
	"github.com/patrickhuber/go-wasm/api"
)

type ModuleInstance struct {
	store             *Store
	Types             []api.FuncType
	FunctionAddresses []address.Function
	Exports           []ExportInstance
}

type ExportInstance struct {
	Name  string
	Value ExternalValue
}

type ExternalValue interface {
	externalValue()
}

type FunctionExternalValue struct {
	Func address.Function
}

func NewModuleInstance(store *Store, module *api.Module) (*ModuleInstance, error) {
	moduleInstance := &ModuleInstance{
		store: store,
	}
	for _, fn := range module.Funcs {
		funcAddr := len(store.Funcs)
		store.Funcs = append(store.Funcs, fn)
		moduleInstance.FunctionAddresses = append(moduleInstance.FunctionAddresses, address.Function(funcAddr))
	}
	return moduleInstance, nil
}
