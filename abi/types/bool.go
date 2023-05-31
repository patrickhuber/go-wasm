package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Bool bool

func (Bool) Kind() kind.Kind {
	return kind.Bool
}

func (Bool) Size() uint32 {
	return 1
}

func (Bool) Alignment() uint32 {
	return 1
}
