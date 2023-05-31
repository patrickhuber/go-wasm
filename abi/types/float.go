package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Float32 float32

func (Float32) Kind() kind.Kind {
	return kind.Float32
}

func (Float32) Size() uint32 {
	return 4
}

func (Float32) Alignment() uint32 {
	return 4
}

type Float64 float64

func (Float64) Kind() kind.Kind {
	return kind.Float64
}

func (Float64) Size() uint32 {
	return 8
}

func (Float64) Alignment() uint32 {
	return 8
}
