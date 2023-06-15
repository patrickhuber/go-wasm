package types

import "github.com/patrickhuber/go-wasm/abi/kind"

const (
	MaxStringByteLength uint32 = (1 << 31) - 1
)

type String struct{}

func (String) Kind() kind.Kind {
	return kind.String
}

func (String) Size() uint32 {
	return 8
}

func (String) Alignment() uint32 {
	return 4
}

func (s String) Despecialize() ValType {
	return s
}

func (String) Flatten() []kind.Kind {
	return []kind.Kind{kind.U32, kind.U32}
}
