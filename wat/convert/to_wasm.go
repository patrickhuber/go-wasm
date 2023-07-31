package convert

import (
	"github.com/patrickhuber/go-wasm/component"
	"github.com/patrickhuber/go-wasm/module"
	"github.com/patrickhuber/go-wasm/types"
	"github.com/patrickhuber/go-wasm/wat/ast"
)

func ToComponent(c *ast.Component) *component.Component {
	return nil
}

func ToModule(m *ast.Module) *module.Module {
	mod := &module.Module{}
	for _, f := range m.Functions {

		mod.Types = append(mod.Types, *ToFuncType(&f))
		mod.Functions = append(mod.Functions, *ToFunction(&f))
	}
	return mod
}

func ToFunction(f *ast.Function) *module.Function {

	return &module.Function{}
}

func ToFuncType(f *ast.Function) *module.FuncType {
	ft := &module.FuncType{}
	for _, p := range f.Parameters {
		param := ToType(p.Type)
		ft.Parameters = append(ft.Parameters, param)
	}
	for _, r := range f.Results {
		result := ToType(r.Type)
		ft.Results = append(ft.Results, result)
	}
	return ft
}

func ToType(f ast.Type) types.Value {
	return nil
}
