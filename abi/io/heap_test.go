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
