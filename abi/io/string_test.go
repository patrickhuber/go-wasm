package io_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/abi/io"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	encodings := []encoding.Encoding{
		encoding.UTF8,
		encoding.UTF16,
		encoding.Latin1Utf16,
	}

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
		"a\ud7ffb",
		"a\u02ff\u03ff\u04ffbc",
		"\uf123",
		"\uf123\uf123abc",
		"abcdef\uf123",
	}

	for i, test := range strings {
		for _, src := range encodings {
			for _, dst := range encodings {
				name := fmt.Sprintf("from: %s to: %s string: %d", src, dst, i)
				t.Run(name, func(t *testing.T) {
					err := testString(src, dst, test)
					require.Nil(t, err)
				})
			}
		}
	}

}

func testString(srcEncoding encoding.Encoding, dstEncoding encoding.Encoding, s string) error {
	fallback := encoding.None
	if srcEncoding == encoding.Latin1Utf16 {
		srcEncoding = encoding.Latin1
		fallback = encoding.UTF16
	}

	sourceCodec, err := encoding.DefaultFactory().Get(srcEncoding)
	if err != nil {
		return err
	}

	encoded, err := encoding.EncodeString(sourceCodec, s)
	if err == nil {
		return testStringInternal(srcEncoding, dstEncoding, s, encoded, len(encoded)/sourceCodec.RuneSize())
	}

	if fallback == encoding.None {
		return err
	}

	srcEncoding = fallback

	sourceCodec, err = encoding.DefaultFactory().Get(srcEncoding)
	if err != nil {
		return err
	}

	encoded, err = encoding.EncodeString(sourceCodec, s)
	if err != nil {
		return err
	}

	tcu := io.TaggedCodeUnits{
		CodeUnits: uint32(len(encoded) / sourceCodec.RuneSize()),
		UTF16:     true,
	}

	codeUnits := int(io.TaggedCodeUnitsToUint32(tcu))
	return testStringInternal(srcEncoding, dstEncoding, s, encoded, codeUnits)
}

func testStringInternal(srcEncoding encoding.Encoding, dstEncoding encoding.Encoding, s string, encoded []byte, taggedCodeUnits int) error {
	heap := NewHeap(len(encoded))
	buf := heap.Memory.Bytes()
	copy(buf, encoded)
	cx := NewContext(heap.Memory, srcEncoding, heap.ReAllocate, nil)
	return test(types.String{}, []any{int32(0), int32(taggedCodeUnits)}, s, cx, dstEncoding, nil, nil)
}
