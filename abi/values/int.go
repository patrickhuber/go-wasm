package values

import "github.com/patrickhuber/go-wasm/abi/kind"

type S32 int32

func (S32) Kind() kind.Kind {
	return kind.S32
}

func (i S32) Value() any {
	return int32(i)
}

type S64 int64

func (S64) Kind() kind.Kind {
	return kind.S64
}

func (i S64) Value() any {
	return int64(i)
}
