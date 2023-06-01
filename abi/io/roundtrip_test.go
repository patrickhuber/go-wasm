package io_test

import (
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

	c := &types.Context{
		Options: &types.CanonicalOptions{
			Memory: make([]byte, 8),
		},
	}
	for i, test := range tests {

		t.Run(test.t.Kind().String(), func(t *testing.T) {
			zero(c.Options.Memory)
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
