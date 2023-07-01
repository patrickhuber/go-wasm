package io_test

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func TestCanRoundTrip(t *testing.T) {
	type test struct {
		val any
		t   types.ValType
	}
	tests := []test{
		{uint8(math.MaxUint8), U8()},
		{uint16(math.MaxUint16), U16()},
		{uint32(math.MaxUint32), U32()},
		{uint64(math.MaxUint64), U64()},
		{int8(math.MaxInt8), S8()},
		{int16(math.MaxInt16), S16()},
		{int32(math.MaxInt32), S32()},
		{int64(math.MaxInt64), S64()},
		{float32(math.MaxFloat32), Float32()},
		{float64(math.MaxFloat64), Float64()},
	}

	heap := NewHeap(8)
	c := NewContext(heap.Memory, encoding.UTF8, heap.ReAllocate, func() {})
	for i, test := range tests {
		name := reflect.ValueOf(test.t).Elem().Type().Name()
		t.Run(name, func(t *testing.T) {
			zero(c.Options.Memory.Bytes())
			err := io.Store(c, test.val, test.t, 0)
			require.Nil(t, err, "store %d type %T", i, test.t)

			val, err := io.Load(c, test.t, 0)
			require.Nil(t, err, "load %d type %T", i, test.t)
			require.Equal(t, test.val, val, "load %d type %T", i, test.t)
		})

	}
}

func zero[T byte](slice []T) {
	for i := 0; i < len(slice); i++ {
		slice[i] = 0
	}
}

type Heap struct {
	Memory    *bytes.Buffer
	LastAlloc int
}

func NewHeap(size int) *Heap {
	return &Heap{
		Memory:    bytes.NewBuffer(make([]byte, size)),
		LastAlloc: 0,
	}
}

func (h *Heap) ReAllocate(originalPtr, originalSize, alignment, newSize uint32) (uint32, error) {
	if originalPtr != 0 && newSize < originalSize {
		return io.AlignTo(originalPtr, alignment)
	}

	ret, err := io.AlignTo(uint32(h.LastAlloc), alignment)
	if err != nil {
		return 0, err
	}
	h.LastAlloc = int(ret + newSize)

	// are we over the capacity?
	if h.LastAlloc > h.Memory.Cap() {
		return 0, fmt.Errorf("Out of Memory: Have %d need %d", h.Memory.Cap(), h.LastAlloc)
	}

	h.Memory.Grow(h.LastAlloc)

	// memcopy here?
	buf := h.Memory.Bytes()
	copy(buf[ret:ret+originalSize], buf[originalPtr:originalPtr+originalSize])

	return ret, nil
}

func NewContext(memory *bytes.Buffer, enc encoding.Encoding, realloc types.ReallocFunc, postReturn types.PostReturnFunc) *types.CallContext {
	options := NewOptions(memory, enc, realloc, postReturn)
	return &types.CallContext{
		Options: options,
	}
}

func NewOptions(memory *bytes.Buffer, enc encoding.Encoding, realloc types.ReallocFunc, postReturn types.PostReturnFunc) *types.CanonicalOptions {
	return &types.CanonicalOptions{
		Memory:         memory,
		StringEncoding: enc,
		Realloc:        realloc,
		PostReturn:     postReturn,
	}
}
