package types

import "github.com/patrickhuber/go-wasm/abi/kind"

const (
	MaxStringByteLength uint32 = (1 << 31) - 1
)

type String struct{}

func (String) Kind() kind.Kind {
	return kind.String
}

func (String) Size() (uint32, error) {
	return 8, nil
}

func (String) Alignment() (uint32, error) {
	return 4, nil
}

func (s String) Despecialize() ValType {
	return s
}

func (String) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32, kind.U32}, nil
}
