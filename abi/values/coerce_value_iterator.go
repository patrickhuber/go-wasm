package values

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

type CoerceValueIterator interface {
	ValueIterator
	ValueIterator() ValueIterator
	FlatTypes() []kind.Kind
}

type coerceValueIterator struct {
	inner     ValueIterator
	flatTypes []kind.Kind
}

func (vi *coerceValueIterator) ValueIterator() ValueIterator {
	return vi.inner
}

func (vi *coerceValueIterator) FlatTypes() []kind.Kind {
	return vi.flatTypes
}

// Index implements ValueIterator.
func (vi *coerceValueIterator) Index() int {
	return vi.inner.Index()
}

// Length implements ValueIterator.
func (vi *coerceValueIterator) Length() int {
	return vi.inner.Length()
}

// Next implements ValueIterator.
func (vi *coerceValueIterator) Next(want kind.Kind) (any, error) {
	if vi.inner == nil {
		return nil, fmt.Errorf("inner value iterator is nil")
	}
	if len(vi.flatTypes) == 0 {
		return nil, fmt.Errorf("inner flat types is nil")
	}
	have := vi.flatTypes[0]
	vi.flatTypes = vi.flatTypes[1:]
	x, err := vi.inner.Next(have)
	if err != nil {
		return nil, err
	}
	switch {
	case have == kind.U32 && want == kind.Float32:
		u32, ok := x.(uint32)
		if !ok {
			return nil, types.NewCastError(x, "uint32")
		}
		return float32(u32), nil
	case have == kind.U64 && want == kind.U32:
		u64, ok := x.(uint64)
		if !ok {
			return nil, types.NewCastError(x, "uint64")
		}
		return uint32(u64), nil
	case have == kind.U64 && want == kind.Float32:
		u64, ok := x.(uint64)
		if !ok {
			return nil, types.NewCastError(x, "uint64")
		}
		return float32(uint32(u64)), nil
	case have == kind.U64 && want == kind.Float64:
		u64, ok := x.(uint64)
		if !ok {
			return nil, types.NewCastError(x, "uint64")
		}
		return float64(u64), nil
	default:
		return x, nil
	}
}

func NewCoerceValueIterator(inner ValueIterator, flatTypes []kind.Kind) CoerceValueIterator {
	return &coerceValueIterator{
		inner:     inner,
		flatTypes: flatTypes,
	}
}
