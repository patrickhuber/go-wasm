package types

import (
	"strconv"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

type Union struct {
	Types []ValType
}

func (*Union) Kind() kind.Kind {
	return kind.Union
}

func (u *Union) Size() (uint32, error) {
	return u.Despecialize().Size()
}

func (u *Union) Alignment() (uint32, error) {
	return u.Despecialize().Alignment()
}

func (u *Union) Despecialize() ValType {
	var cases []Case
	for i, t := range u.Types {
		c := Case{
			Label: strconv.Itoa(i),
			Type:  t,
		}
		cases = append(cases, c)
	}
	return &Variant{
		Cases: cases,
	}
}
