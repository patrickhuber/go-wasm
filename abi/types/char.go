package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Char rune

func (Char) Kind() kind.Kind {
	return kind.Char
}

func (Char) Size() (uint32, error) {
	return 4, nil
}

func (Char) Alignment() (uint32, error) {
	return 4, nil
}

func (c Char) Despecialize() ValType {
	return c
}

func (Char) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}
