package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Char rune

func (Char) Kind() kind.Kind {
	return kind.Char
}

func (Char) Size() uint32 {
	return 4
}

func (Char) Alignment() uint32 {
	return 4
}

func (c Char) Despecialize() ValType {
	return c
}

func (Char) Flatten() []kind.Kind {
	return []kind.Kind{kind.S32}
}
