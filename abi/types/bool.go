package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Bool struct{}

func (Bool) Kind() kind.Kind {
	return kind.Bool
}

func (Bool) Size() uint32 {
	return 1
}

func (Bool) Alignment() uint32 {
	return 1
}

func (b Bool) Despecialize() ValType {
	return b
}

func (b Bool) Flatten() []kind.Kind {
	return []kind.Kind{kind.S32}
}
