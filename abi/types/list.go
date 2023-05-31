package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type List struct {
	Type ValType
}

func (*List) Kind() kind.Kind {
	return kind.List
}

func (*List) Size() uint32 {
	return 8
}

func (*List) Alignment() uint32 {
	return 4
}

func (l *List) Despecialize() ValType {
	return l
}
