package runtime

import (
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/wasm"
)

type Instantiator interface {
	Instantiate(io.Reader) (*ObjectInstance, error)
}

type instantiator struct {
}

func (i *instantiator) Instantiate(reader io.Reader) (*ObjectInstance, error) {
	objReader := wasm.NewObjectDecoder(reader)

	obj, err := objReader.Decode()
	if err != nil {
		return nil, err
	}

	return i.instantiateObject(obj)
}

func (i *instantiator) instantiateObject(obj *wasm.Object) (*ObjectInstance, error) {
	var component *ComponentInstance
	var module *ModuleInstance
	var err error

	switch {
	case obj.Component == nil:
		component, err = i.instantiateComponent(obj.Component)
		if err != nil {
			return nil, err
		}
	case obj.Module == nil:
		module, err = i.instantiateModule(obj.Module)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid object. Module or Component required")
	}
	return &ObjectInstance{
		Module:    module,
		Component: component,
	}, nil
}

func (i *instantiator) instantiateModule(module *wasm.Module) (*ModuleInstance, error) {
	instance := &ModuleInstance{}
	for _, section := range module.Sections {
		switch section.ID {
		case wasm.CustomSectionType:
		case wasm.TypeSectionType:
			t := section.Type
			if t == nil {
				return nil, fmt.Errorf("expected type in section to not be nil")
			}
		case wasm.ImportSectionType:
		case wasm.FuncSectionType:
			fun := section.Function
			if fun == nil {
				return nil, fmt.Errorf("expected function in section to not be nil")
			}

		case wasm.TableSectionType:
		case wasm.MemSectionType:
		case wasm.GlobalSectionType:
		case wasm.ExportSectionType:
		case wasm.StartSectionType:
		case wasm.ElemSectionType:
		case wasm.CodeSectionType:
		case wasm.DataSectionType:
		default:
			return nil, fmt.Errorf("invalid section type")
		}
	}
	return instance, nil
}

func (i *instantiator) instantiateComponent(component *wasm.Component) (*ComponentInstance, error) {
	return nil, nil
}
