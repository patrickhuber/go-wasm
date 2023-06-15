package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Float32 struct{}

func (Float32) Kind() kind.Kind {
	return kind.Float32
}

func (Float32) Size() (uint32, error) {
	return 4, nil
}

func (Float32) Alignment() (uint32, error) {
	return 4, nil
}

func (f32 Float32) Despecialize() ValType {
	return f32
}

func (f32 Float32) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.Float32}, nil
}

type Float64 struct{}

func (Float64) Kind() kind.Kind {
	return kind.Float64
}

func (Float64) Size() (uint32, error) {
	return 8, nil
}

func (Float64) Alignment() (uint32, error) {
	return 8, nil
}

func (f64 Float64) Despecialize() ValType {
	return f64
}

func (f64 Float64) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.Float64}, nil
}
