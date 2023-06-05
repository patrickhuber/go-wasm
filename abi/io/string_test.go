package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	strings := []string{
		"",
		"a",
		"hi",
		"\x00",
		"a\x00b",
		"\x00b",
		"\x80",
		"\x80b",
		"ab\xefc",
		"\u01ffy",
		"xy\u01ff",

		"abcdef\uf123",
	}

	for i, test := range strings {
		name := fmt.Sprintf("iteration: %d", i)
		t.Run(name, func(t *testing.T) {
			err := testString(types.Utf8, types.Utf8, test)
			require.Nil(t, err)
		})
	}
}

func testString(srcEncoding types.StringEncoding, dstEncoding types.StringEncoding, s string) error {
	switch srcEncoding {
	case types.Utf8:
		return testStringInternal(types.Utf8, dstEncoding, s, s, len(s))
	default:
		return fmt.Errorf("invalid source encoding '%v'", srcEncoding)
	}
}

func testStringInternal(srcEncoding types.StringEncoding, dstEncoding types.StringEncoding, s string, encoded string, taggedCodeUnits int) error {
	heap := NewHeap(len(encoded))
	buf := heap.Memory.Bytes()
	copy(buf, []byte(encoded))
	cx := NewContext(heap.Memory, srcEncoding, heap.ReAllocate, nil)
	return test(types.String{}, []any{int32(0), int32(taggedCodeUnits)}, s, cx, dstEncoding, nil, nil)
}
