package values

import "github.com/patrickhuber/go-wasm/abi/kind"

type Float32 float32

func (Float32) Kind() kind.Kind {
	return kind.Float32
}

func (i Float32) Value() any {
	return float32(i)
}

type Float64 float64

func (Float64) Kind() kind.Kind {
	return kind.Float64
}

func (i Float64) Value() any {
	return float64(i)
}
