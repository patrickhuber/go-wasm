package values

import (
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

type Value interface {
	Kind() kind.Kind
	Value() any
}

type ValueIterator interface {
	Next(k kind.Kind) (any, error)
}

type valueIterator struct {
	values []Value
	index  int
}

func (vi *valueIterator) Next(k kind.Kind) (any, error) {
	v := vi.values[vi.index]
	vi.index += 1
	if v.Kind() != k {
		return nil, types.Trap()
	}
	return v.Value(), nil
}
