package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	// hex literals will fail because they are not converted to utf8
	// to work around this, use unicode literals instead
	strings := []string{
		"",
		"a",
		"hi",
		"\u0000",    // "\x00",
		"a\u0000b",  // "a\x00b",
		"\u0000b",   // "\x00b",
		"\u0080",    // "\x80",
		"\u0080b",   // "\x80b",
		"ab\u00efc", // "ab\xefc",
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
		return testStringInternal(srcEncoding, dstEncoding, s, []byte(s), len(s))
	default:
		return fmt.Errorf("invalid source encoding '%v'", srcEncoding)
	}
}

func testStringInternal(srcEncoding types.StringEncoding, dstEncoding types.StringEncoding, s string, encoded []byte, taggedCodeUnits int) error {
	heap := NewHeap(len(encoded))
	buf := heap.Memory.Bytes()
	copy(buf, encoded)
	cx := NewContext(heap.Memory, srcEncoding, heap.ReAllocate, nil)
	return test(types.String{}, []any{int32(0), int32(taggedCodeUnits)}, s, cx, dstEncoding, nil, nil)
}
