package io_test

import (
	"reflect"
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
	testPairs[uint32, uint8](t, U8(),
		Pair(uint32(127), uint8(127)),
		Pair(uint32(128), uint8(128)),
		Pair(uint32(255), uint8(255)),
		Pair(uint32(256), uint8(0)),
		Pair(uint32(4294967295), uint8(255)),
		Pair(uint32(4294967168), uint8(128)),
		Pair(uint32(4294967167), uint8(127)))
	testPairs[uint32, int8](t, S8(),
		Pair(uint32(127), int8(127)),
		Pair(uint32(128), int8(-128)),
		Pair(uint32(255), int8(-1)),
		Pair(uint32(256), int8(0)),
		Pair(uint32(4294967295), int8(-1)),
		Pair(uint32(4294967168), int8(-128)),
		Pair(uint32(4294967167), int8(127)))
	testPairs[uint32, uint16](t, U16(),
		Pair(uint32(32767), uint16(32767)),
		Pair(uint32(32768), uint16(32768)),
		Pair(uint32(65535), uint16(65535)),
		Pair(uint32(65536), uint16(0)),
		Pair(uint32((1<<32)-1), uint16(65535)),
		Pair(uint32((1<<32)-32768), uint16(32768)),
		Pair(uint32((1<<32)-32769), uint16(32767)))
	testPairs[uint32, int16](t, S16(),
		Pair(uint32(32767), int16(32767)),
		Pair(uint32(32768), int16(-32768)),
		Pair(uint32(65535), int16(-1)),
		Pair(uint32(65536), int16(0)),
		Pair(uint32((1<<32)-1), int16(-1)),
		Pair(uint32((1<<32)-32768), int16(-32768)),
		Pair(uint32((1<<32)-32769), int16(32767)))
	testPairs[uint32, uint32](t, U32(),
		Pair(uint32((1<<31)-1), uint32(1<<31)-1),
		Pair(uint32(1<<31), uint32(1<<31)),
		Pair(uint32((1<<32)-1), uint32((1<<32)-1)))
	testPairs[uint32, int32](t, S32(),
		Pair(uint32((1<<31)-1), int32((1<<31)-1)),
		Pair(uint32(1<<31), int32(-(1<<31))),
		Pair(uint32((1<<32)-1), int32(-1)))
	testPairs[uint64, uint64](t, U64(),
		Pair(uint64((1<<63)-1), uint64((1<<63)-1)),
		Pair(uint64(1<<63), uint64(1<<63)),
		Pair(uint64((1<<64)-1), uint64((1<<64)-1)))
	testPairs[uint64, int64](t, S64(),
		Pair(uint64((1<<63)-1), int64((1<<63)-1)),
		Pair(uint64(1<<63), int64(-(1<<63))),
		Pair(uint64((1<<64)-1), int64(-1)))
	testPairs[float32, float32](t, Float32(),
		Pair(float32(3.14), float32(3.14)))
	testPairs[float64, float64](t, Float64(),
		Pair(float64(3.14), float64(3.14)))
	testPairs[uint32, rune](t, Char(),
		Pair(uint32(0), '\x00'),
		Pair(uint32(65), 'A'),
		Pair(uint32(0xD7FF), '\uD7FF'),
		Pair(uint32(0xE000), '\uE000'),
		Pair(uint32(0x10FFFF), '\U0010FFFF'))
	testPairs[uint32, any](t, Char(), // any nil values must be passed as 'any'
		Pair[uint32, any](uint32(0xD800), nil),
		Pair[uint32, any](uint32(0xDFFF), nil),
		Pair[uint32, any](uint32(0x110000), nil),
		Pair[uint32, any](uint32(0xFFFFFFFF), nil))
	testPairs[uint32, map[string]any](t, Enum("a", "b"),
		Pair(uint32(0), map[string]any{"a": nil}),
		Pair(uint32(1), map[string]any{"b": nil}))
	testPairs[uint32, any](t, Enum("a", "b"),
		Pair[uint32, any](uint32(2), nil)) // any nil values must be passed as 'any'
}
func testPairs[TLift, TValue any](t *testing.T, vt types.ValType, pairs ...pair[TLift, TValue]) {
	for _, p := range pairs {
		name := reflect.ValueOf(vt).Elem().Type().Name()
		t.Run(name, func(t *testing.T) {
			cxt := Context()
			err := test(vt, []any{p.ValsToLift}, p.Value, cxt, encoding.UTF8, nil, nil)
			require.Nil(t, err)
		})
	}
}
