package encoding_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if !ok {
			return
		}
		t.Fatal(err)
	}()

	encodings := []encoding.Encoding{
		encoding.UTF8,
		encoding.UTF16,
		encoding.UTF16BE,
		encoding.UTF16LE,
		encoding.Latin1,
	}

	// some hex literals will fail because they are not converted to utf8
	// to work around this, use unicode literals instead
	tests := []string{
		"",
		"a",
		"hi",
		"\x00",
		"a\x00b",
		"\x00b",
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

	factory := encoding.DefaultFactory()
	fallback := factory.Get(encoding.UTF16LE).Unwrap()

	for i, test := range tests {
		for _, enc := range encodings {

			codec := factory.Get(enc).Unwrap()
			name := fmt.Sprintf("iteration %d codec %s", i, codec.Encoding())
			t.Run(name, func(t *testing.T) {

				buf, err := encoding.EncodeString(codec, test)
				// retry for latin1
				if err != nil && codec.Encoding() == encoding.Latin1 {
					codec = fallback
					buf, err = encoding.EncodeString(codec, test)
				}
				require.Nil(t, err)

				str, err := encoding.DecodeString(codec, bytes.NewReader(buf))
				require.Nil(t, err)

				require.Equal(t, test, str)
			})
		}
	}
}
