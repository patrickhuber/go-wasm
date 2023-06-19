package io_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

type pair[TLift, TValue any] struct {
	ValsToLift TLift
	Value      TValue
}

func Pair[TLift, TValue any](lift TLift, value TValue) pair[TLift, TValue] {
	return pair[TLift, TValue]{Value: value, ValsToLift: lift}
}

func TestPairs(t *testing.T) {
	testPairs[uint32, bool](t, Bool(),
		Pair(uint32(0), false),
		Pair(uint32(1), true),
		Pair(uint32(2), true),
		Pair(uint32(4294967295), true))
}

func testPairs[TLift, TValue any](t *testing.T, vt types.ValType, pairs ...pair[TLift, TValue]) {
	for _, p := range pairs {
		t.Run(vt.Kind().String(), func(t *testing.T) {
			cxt := NewContext(&bytes.Buffer{}, encoding.UTF8, nil, nil)
			err := test(vt, []any{p.ValsToLift}, p.Value, cxt, encoding.UTF8, nil, nil)
			require.Nil(t, err)
		})
	}
}
