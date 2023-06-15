package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type List struct {
	Type ValType
}

func (List) Kind() kind.Kind {
	return kind.List
}

func (List) Size() (uint32, error) {
	return 8, nil
}

func (List) Alignment() (uint32, error) {
	return 4, nil
}

func (l *List) Despecialize() ValType {
	return l
}

func (List) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32, kind.U32}, nil
}
