package io_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func TestHeap(t *testing.T) {
	type test struct {
		name     string
		vt       types.ValType
		expected any
		args     []any
		bytes    []byte
	}
	tests := []test{
		{"list_record", List(Record()), []any{map[string]any{}, map[string]any{}, map[string]any{}}, []any{uint32(0), uint32(3)}, []byte{}},
		{"list_bool", List(Bool()), []any{true, false, true}, []any{uint32(0), uint32(3)}, []byte{1, 0, 1}},
		{"list_bool", List(Bool()), []any{true, false, true}, []any{uint32(0), uint32(3)}, []byte{1, 0, 2}},
		{"list_bool", List(Bool()), []any{true, false, true}, []any{uint32(3), uint32(3)}, []byte{0xff, 0xff, 0xff, 1, 0, 1}},
		{"list_u8", List(U8()), []any{uint8(1), uint8(2), uint8(3)}, []any{uint32(0), uint32(3)}, []byte{1, 2, 3}},
		{"list_u16", List(U16()), []any{uint16(1), uint16(2), uint16(3)}, []any{uint32(0), uint32(3)}, []byte{1, 0, 2, 0, 3, 0}},
		{"list_u16", List(U16()), nil, []any{uint32(1), uint32(3)}, []byte{0, 1, 0, 2, 0, 3, 0}},
		{"list_u32", List(U32()), []any{uint32(1), uint32(2), uint32(3)}, []any{uint32(0), uint32(3)}, []byte{1, 0, 0, 0, 2, 0, 0, 0, 3, 0, 0, 0}},
		{"list_u32", List(U64()), []any{uint64(1), uint64(2)}, []any{uint32(0), uint32(2)}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0}},
		{"list_s8", List(S8()), []any{int8(-1), int8(-2), int8(-3)}, []any{uint32(0), uint32(3)}, []byte{0xff, 0xfe, 0xfd}},
		{"list_s16", List(S16()), []any{int16(-1), int16(-2), int16(-3)}, []any{uint32(0), uint32(3)}, []byte{0xff, 0xff, 0xfe, 0xff, 0xfd, 0xff}},
		{"list_s32", List(S32()), []any{int32(-1), int32(-2), int32(-3)}, []any{uint32(0), uint32(3)}, []byte{0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xfd, 0xff, 0xff, 0xff}},
		{"list_s64", List(S64()), []any{int64(-1), int64(-2)}, []any{uint32(0), uint32(2)}, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"list_char", List(Char()), []any{'A', 'B', 'c'}, []any{uint32(0), uint32(3)}, []byte{65, 00, 00, 00, 66, 00, 00, 00, 99, 00, 00, 00}},
		{"list_string", List(String()), []any{"hi", "wat"}, []any{uint32(0), uint32(2)}, []byte{16, 0, 0, 0, 2, 0, 0, 0, 21, 0, 0, 0, 3, 0, 0, 0,
			uint8('h'), uint8('i'), 0xf, 0xf, 0xf, uint8('w'), uint8('a'), uint8('t')}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHeap(t, test.vt, test.expected, test.args, test.bytes)
		})
	}
}

func testHeap(t *testing.T, vt types.ValType, expect any, args []any, bytes []byte) {
	heap := NewHeap(len(bytes))
	copy(heap.Memory.Bytes(), bytes)

	cx := NewContext(heap.Memory, encoding.UTF8, nil, nil)
	err := test(vt, args, expect, cx, cx.Options.StringEncoding, nil, nil)

	require.Nil(t, err)
}
