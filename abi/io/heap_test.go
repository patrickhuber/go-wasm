package io_test

import (
	"strconv"
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
			byte('h'), byte('i'), 0xf, 0xf, 0xf, byte('w'), byte('a'), byte('t')}},
		{"list_list_u8", List(List(U8())), []any{[]any{byte(3), byte(4), byte(5)}, []any(nil), []any{byte(6), byte(7)}}, []any{uint32(0), uint32(3)}, []byte{24, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 27, 0, 0, 0, 2, 0, 0, 0, 3, 4, 5, 6, 7}},
		{"list_list_u16", List(List(U16())), []any{[]any{uint16(5), uint16(6)}}, []any{uint32(0), uint32(1)}, []byte{8, 0, 0, 0, 2, 0, 0, 0, 5, 0, 6, 0}},
		{"list_list_u16", List(List(U16())), nil, []any{uint32(0), uint32(1)}, []byte{9, 0, 0, 0, 2, 0, 0, 0, 0, 5, 0, 6, 0}},
		{"list_tuple_u8_u8_u16_u32", List(Tuple(U8(), U8(), U16(), U32())), []any{NewTuple(byte(6), byte(7), uint16(8), uint32(9)), NewTuple(byte(4), byte(5), uint16(6), uint32(7))}, []any{uint32(0), uint32(2)}, []byte{6, 7, 8, 0, 9, 0, 0, 0, 4, 5, 6, 0, 7, 0, 0, 0}},
		{"list_tuple_u8_u16_u8_u32", List(Tuple(U8(), U16(), U8(), U32())), []any{NewTuple(byte(6), uint16(7), byte(8), uint32(9)), NewTuple(byte(4), uint16(5), byte(6), uint32(7))}, []any{uint32(0), uint32(2)}, []byte{6, 0xff, 7, 0, 8, 0xff, 0xff, 0xff, 9, 0, 0, 0, 4, 0xff, 5, 0, 6, 0xff, 0xff, 0xff, 7, 0, 0, 0}},
		{"list_tuple_u16_u8", List(Tuple(U16(), U8())), []any{NewTuple(uint16(6), uint8(7)), NewTuple(uint16(8), uint8(9))}, []any{uint32(0), uint32(2)}, []byte{6, 0, 7, 0x0ff, 8, 0, 9, 0xff}},
		{"list_tuple_tuple_u16_u8_u8", List(Tuple(Tuple(U16(), U8()), U8())), []any{NewTuple(NewTuple(uint16(4), uint8(5)), uint8(6)), NewTuple(NewTuple(uint16(7), uint8(8)), uint8(9))}, []any{uint32(0), uint32(2)}, []byte{4, 0, 5, 0xff, 6, 0xff, 7, 0, 8, 0xff, 9, 0xff}},
		{"list_union_record_u8_tuple_u8_u16", List(Union(Record(), U8(), Tuple(U8(), U16()))), []any{map[string]any{"0": map[string]any{}}, map[string]any{"1": byte(42)}, map[string]any{"2": NewTuple(byte(6), uint16(7))}}, []any{uint32(0), uint32(3)}, []byte{0, 0xff, 0xff, 0xff, 0xff, 0xff, 1, 0xff, 42, 0xff, 0xff, 0xff, 2, 0xff, 6, 0xff, 7, 0}},
		{"list_union_u32_u8", List(Union(U32(), U8())), []any{map[string]any{"0": uint32(256)}, map[string]any{"1": uint8(42)}}, []any{uint32(0), uint32(2)}, []byte{0, 0xff, 0xff, 0xff, 0, 1, 0, 0, 1, 0xff, 0xff, 0xff, 42, 0xff, 0xff, 0xff}},
		{"list_tuple_union_u8_tuple_u16_u8_u8", List(Tuple(Union(U8(), Tuple(U16(), U8())), U8())), []any{NewTuple(map[string]any{"1": NewTuple(uint16(5), uint8(6))}, uint8(7)), NewTuple(map[string]any{"0": uint8(8)}, uint8(9))}, []any{uint32(0), uint32(2)}, []byte{1, 0xff, 5, 0, 6, 0xff, 7, 0xff, 0, 0xff, 8, 0xff, 0xff, 0xff, 9, 0xff}},
		{"list_union_u8", List(Union(U8())), []any{map[string]any{"0": uint8(6)}, map[string]any{"0": uint8(7)}, map[string]any{"0": uint8(8)}}, []any{uint32(0), uint32(3)}, []byte{0, 6, 0, 7, 0, 8}},
		{"list_flags", List(Flags()), []any{map[string]any{}, map[string]any{}, map[string]any{}}, []any{uint32(0), uint32(3)}, []byte{}},
		{"list_tuple_flags_u8", List(Tuple(Flags(), U8())), []any{NewTuple(map[string]any{}, uint8(42)), NewTuple(map[string]any{}, uint8(43)), NewTuple(map[string]any{}, uint8(44))}, []any{uint32(0), uint32(3)}, []byte{42, 43, 44}},
		{"list_flags", List(Flags("a", "b")), []any{map[string]any{"a": false, "b": false}, map[string]any{"a": false, "b": true}, map[string]any{"a": true, "b": true}}, []any{uint32(0), uint32(3)}, []byte{0, 2, 3}},
		{"list_flags", List(Flags("a", "b")), []any{map[string]any{"a": false, "b": false}, map[string]any{"a": false, "b": true}, map[string]any{"a": false, "b": false}}, []any{uint32(0), uint32(3)}, []byte{0, 2, 4}},
		// {"", nil, []any{}, []any{}, []byte{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testHeap(t, test.vt, test.expected, test.args, test.bytes)
		})
	}
}

func NewTuple(values ...any) map[string]any {
	m := map[string]any{}
	for i, value := range values {
		m[strconv.Itoa(i)] = value
	}
	return m
}

func testHeap(t *testing.T, vt types.ValType, expect any, args []any, bytes []byte) {
	heap := NewHeap(len(bytes))
	copy(heap.Memory.Bytes(), bytes)

	cx := NewContext(heap.Memory, encoding.UTF8, nil, nil)
	err := test(vt, args, expect, cx, cx.Options.StringEncoding, nil, nil)

	require.Nil(t, err)
}
