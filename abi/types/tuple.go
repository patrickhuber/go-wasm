package types

import (
	"strconv"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

type Tuple struct {
	Types []ValType
}

func (t *Tuple) Kind() kind.Kind {
	return kind.Tuple
}

func (t *Tuple) Size() uint32 {
	return t.Despecialize().Size()
}

func (t *Tuple) Alignment() uint32 {
	return t.Despecialize().Alignment()
}

func (t *Tuple) Despecialize() ValType {
	var fields []Field
	for i, ty := range t.Types {
		field := Field{
			Label: strconv.Itoa(i),
			Type:  ty,
		}
		fields = append(fields, field)
	}
	return &Record{
		Fields: fields,
	}
}
