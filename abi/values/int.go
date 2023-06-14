package values

import "github.com/patrickhuber/go-wasm/abi/kind"

type U32 uint32

func (U32) Kind() kind.Kind {
	return kind.U32
}

func (i U32) Value() any {
	return uint32(i)
}

type U64 int64

func (U64) Kind() kind.Kind {
	return kind.U64
}

func (i U64) Value() any {
	return uint64(i)
}
