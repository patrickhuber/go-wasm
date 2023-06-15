package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Bool struct{}

func (Bool) Kind() kind.Kind {
	return kind.Bool
}

func (Bool) Size() (uint32, error) {
	return 1, nil
}

func (Bool) Alignment() (uint32, error) {
	return 1, nil
}

func (b Bool) Despecialize() ValType {
	return b
}

func (b Bool) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}
