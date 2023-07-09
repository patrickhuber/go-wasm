package values

import (
	"math"

	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/assert"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/abi/kind"
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
func (vi *coerceValueIterator) Next(want kind.Kind) (res Result[any]) {
	defer handle.Error(&res)
	assert.NotNilf(vi.inner, "inner value iterator is nil")
	assert.Falsef(len(vi.flatTypes) == 0, "inner flat types is nil")

	have := vi.flatTypes[0]
	vi.flatTypes = vi.flatTypes[1:]
	
	x := vi.inner.Next(have).Unwrap()

	switch {
	case have == kind.U32 && want == kind.Float32:
		u32 := Cast[any, uint32](x).Unwrap()
		return result.Ok[any](math.Float32frombits(u32))
	case have == kind.U64 && want == kind.U32:
		u64 := Cast[any, uint64](x).Unwrap()
		return result.Ok[any](uint32(u64))
	case have == kind.U64 && want == kind.Float32:
		u64 := Cast[any, uint64](x).Unwrap()
		return result.Ok[any](math.Float32frombits(uint32(u64)))
	case have == kind.U64 && want == kind.Float64:
		u64 := Cast[any, uint64](x).Unwrap()
		return result.Ok[any](math.Float64frombits(u64))
	default:
		return result.Ok[any](x)
	}
}

func NewCoerceValueIterator(inner ValueIterator, flatTypes []kind.Kind) CoerceValueIterator {
	return &coerceValueIterator{
		inner:     inner,
		flatTypes: flatTypes,
	}
}
