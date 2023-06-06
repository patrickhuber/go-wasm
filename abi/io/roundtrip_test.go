package io_test

import (
	"bytes"
	"fmt"
	"math"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/stretchr/testify/require"
)

func TestCanRoundTrip(t *testing.T) {
	type test struct {
		val any
		t   types.ValType
	}
	tests := []test{
		{uint8(math.MaxUint8), types.U8{}},
		{uint16(math.MaxUint16), types.U16{}},
		{uint32(math.MaxUint32), types.U32{}},
		{uint64(math.MaxUint64), types.U64{}},
		{int8(math.MaxInt8), types.S8{}},
		{int16(math.MaxInt16), types.S16{}},
		{int32(math.MaxInt32), types.S32{}},
		{int64(math.MaxInt64), types.S64{}},
		{float32(math.MaxFloat32), types.Float32{}},
		{float64(math.MaxFloat64), types.Float64{}},
	}

	heap := NewHeap(8)
	c := NewContext(heap.Memory, types.Utf8, heap.ReAllocate, func() {})
	for i, test := range tests {

		t.Run(test.t.Kind().String(), func(t *testing.T) {
			zero(c.Options.Memory.Bytes())
			err := io.Store(c, test.val, test.t, 0)
			require.Nil(t, err, "store %d type %s", i, test.t.Kind())

			val, err := io.Load(c, test.t, 0)
			require.Nil(t, err, "load %d type %s", i, test.t.Kind())
			require.Equal(t, test.val, val, "load %d type %s", i, test.t.Kind())
		})

	}
}

func zero[T byte](slice []T) {
	for i := 0; i < len(slice); i++ {
		slice[i] = 0
	}
}

type Heap struct {
	Memory    bytes.Buffer
	LastAlloc int
}

func NewHeap(size int) *Heap {
	return &Heap{
		Memory:    *bytes.NewBuffer(make([]byte, size)),
		LastAlloc: 0,
	}
}

func (h *Heap) ReAllocate(originalPtr, originalSize, alignment, newSize uint32) (uint32, error) {
	if originalPtr != 0 && newSize < originalSize {
		return types.AlignTo(originalPtr, alignment), nil
	}

	ret := types.AlignTo(uint32(h.LastAlloc), alignment)
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

func NewContext(memory bytes.Buffer, encoding types.StringEncoding, realloc types.ReallocFunc, postReturn types.PostReturnFunc) *types.Context {
	options := NewOptions(memory, encoding, realloc, postReturn)
	return &types.Context{
		Options: options,
	}
}

func NewOptions(memory bytes.Buffer, encoding types.StringEncoding, realloc types.ReallocFunc, postReturn types.PostReturnFunc) *types.CanonicalOptions {
	return &types.CanonicalOptions{
		Memory:         memory,
		StringEncoding: encoding,
		Realloc:        realloc,
		PostReturn:     postReturn,
	}
}
