package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Enum struct {
	Labels []string
}

func (*Enum) Kind() kind.Kind {
	return kind.Enum
}

func (e *Enum) Size() uint32 {
	vt := e.Despecialize()
	return vt.Size()
}

func (e *Enum) Alignment() uint32 {
	vt := e.Despecialize()
	return vt.Alignment()
}

func (e *Enum) Despecialize() ValType {
	var cases []Case
	for _, v := range e.Labels {
		c := Case{
			Label: v,
			Type:  nil,
		}
		cases = append(cases, c)
	}
	return &Variant{
		Cases: cases,
	}
}

func (e *Enum) Flatten() []kind.Kind {
	return e.Despecialize().Flatten()
}
