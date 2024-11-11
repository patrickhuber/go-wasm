package binary

import (
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/api"
)

type Decoder interface {
	Decode(io.Reader) (api.Directive, error)
}

type decoder struct {
	reader io.Reader
}

func NewDecoder(reader io.Reader) Decoder {
	return &decoder{
		reader: reader,
	}
}

func (decoder *decoder) Decode(reader io.Reader) (api.Directive, error) {
	doc, err := Read(reader)
	if err != nil {
		return nil, err
	}
	return decoder.directive(doc.Directive)
}

func (decoder *decoder) directive(directive Directive) (api.Directive, error) {
	switch d := directive.(type) {
	case *Module:
		return decoder.module(d)
	case *Component:
		return decoder.component(d)
	}
	return nil, fmt.Errorf("unrecognized type %T", directive)
}

func (decoder *decoder) module(module *Module) (*api.Module, error) {
	result := &api.Module{}
	for _, section := range module.Sections {
		switch s := section.(type) {
		case *TypeSection:
			ty, err := decoder.funcTypes(s)
			if err != nil {
				return nil, err
			}
			result.Types = append(result.Types, ty...)
		case *FunctionSection:
			f, err := decoder.function(s)
			if err != nil {
				return nil, err
			}
			result.Funcs = append(result.Funcs, f)
		case *CodeSection:

		}
	}
	return result, nil
}

func (decoder *decoder) component(component *Component) (*api.Component, error) {
	return nil, nil
}

func (decoder *decoder) funcTypes(typeSection *TypeSection) ([]*api.FuncType, error) {
	if typeSection.ID != TypeSectionID {
		return nil, fmt.Errorf("invalid section id. expected %d found %d", TypeSectionID, typeSection.ID)
	}
	var results []*api.FuncType
	for _, ft := range typeSection.Types {
		apift, err := decoder.funcType(ft)
		if err != nil {
			return nil, err
		}
		results = append(results, apift)
	}
	return results, nil
}

func (decoder *decoder) funcType(ft *FunctionType) (*api.FuncType, error) {
	parameters, err := decoder.resultType(&ft.Parameters)
	if err != nil {
		return nil, err
	}
	returns, err := decoder.resultType(&ft.Returns)
	if err != nil {
		return nil, err
	}
	return &api.FuncType{
		Parameters: parameters,
		Returns:    returns,
	}, nil
}

func (decoder *decoder) resultType(rt *ResultType) (api.ResultType, error) {
	return api.ResultType{}, nil
}

func (decoder *decoder) function(funcSection *FunctionSection) (*api.Func, error) {
	return nil, nil
}
