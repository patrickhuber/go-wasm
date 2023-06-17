package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Result struct {
	OK    ValType
	Error ValType
}

func (r *Result) Kind() kind.Kind {
	return kind.Result
}

func (r *Result) Size() (uint32, error) {
	return r.Despecialize().Size()
}

func (r *Result) Alignment() (uint32, error) {
	return r.Despecialize().Alignment()
}

func (r *Result) Despecialize() ValType {
	cases := []Case{
		{
			Label: "ok",
			Type:  r.OK,
		},
		{
			Label: "error",
			Type:  r.Error,
		},
	}
	return &Variant{
		Cases: cases,
	}
}

func (r *Result) Flatten() ([]kind.Kind, error) {
	return r.Despecialize().Flatten()
}
