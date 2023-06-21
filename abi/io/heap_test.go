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
		{"list_record", List(Record()), []any{nil, nil, nil}, []any{uint32(0), uint32(3)}, []byte{}},
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
