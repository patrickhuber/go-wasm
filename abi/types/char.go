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
