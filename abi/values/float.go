package values

import "github.com/patrickhuber/go-wasm/abi/kind"

type Float32Value struct {
	value float32
}

func (i *Float32Value) Kind() kind.Kind {
	return kind.Float64
}

func (i *Float32Value) Value() any {
	return i.value
}

type Float64Value struct {
	value float32
}

func (i *Float64Value) Kind() kind.Kind {
	return kind.Float64
}

func (i *Float64Value) Value() any {
	return i.value
}
