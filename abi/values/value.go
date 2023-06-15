package values

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

type Value interface {
	Kind() kind.Kind
	Value() any
}

type ValueIterator interface {
	Next(k kind.Kind) (any, error)
	Index() int
	Length() int
}

func NewIterator(values ...Value) ValueIterator {
	return &valueIterator{
		index:  0,
		values: values,
	}
}

type valueIterator struct {
	values []Value
	index  int
}

func (vi *valueIterator) Next(k kind.Kind) (any, error) {
	if vi.Length() == 0 {
		return nil, fmt.Errorf("eof")
	}

	v := vi.values[vi.index]
	vi.index += 1
	if v.Kind() != k {
		return nil, fmt.Errorf("error fetching next: have kind.%s, want kind.%s", v.Kind(), k)
	}
	return v.Value(), nil
}

func (vi *valueIterator) Index() int {
	return vi.index
}

func (vi *valueIterator) Length() int {
	return len(vi.values)
}
