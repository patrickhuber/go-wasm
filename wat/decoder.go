package wat

import (
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/wat/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"
	"github.com/patrickhuber/go-wasm/wat/parse"
)

type Decoder interface {
	Decode(io.Reader) (api.Directive, error)
}

func NewDecoder(reader io.Reader) Decoder {
	return &decoder{
		reader: reader,
	}
}

type decoder struct {
	reader io.Reader
}

func (d *decoder) Decode(reader io.Reader) (api.Directive, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	lexer := lex.New(string(buf))
	directive, err := parse.Parse(lexer)
	if err != nil {
		return nil, err
	}
	return d.directive(directive)
}

func (decoder *decoder) directive(directive ast.Directive) (api.Directive, error) {
	switch d := directive.(type) {
	case *ast.Component:
		return decoder.component(d)
	case *ast.Module:
		return decoder.module(d)
	}
	return nil, fmt.Errorf("unrecognized type %T", directive)
}

func (decoder *decoder) component(component *ast.Component) (*api.Component, error) {
	return nil, nil
}

func (decoder *decoder) module(module *ast.Module) (*api.Module, error) {
	return &api.Module{}, nil
}
