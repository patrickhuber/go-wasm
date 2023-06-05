package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Option struct {
	Value ValType
}

func (*Option) Kind() kind.Kind {
	return kind.Option
}

func (o *Option) Size() uint32 {
	return o.Despecialize().Size()
}

func (o *Option) Alignment() uint32 {
	return o.Despecialize().Alignment()
}

func (o *Option) Despecialize() ValType {
	cases := []Case{
		{Label: "none", Type: nil},
		{Label: "some", Type: o.Value},
	}
	return &Variant{
		Cases: cases,
	}
}

func (o *Option) Flatten() []kind.Kind {
	return o.Despecialize().Flatten()
}
