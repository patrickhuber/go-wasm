package encoding_test

import (
	"fmt"
	"testing"

	"github.com/patrickhuber/go-wasm/encoding"
	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	codecs := []encoding.Codec{
		encoding.UTF8(),
		encoding.UTF16(),
	}
	// hex literals will fail because they are not converted to utf8
	// to work around this, use unicode literals instead
	tests := []string{
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

	for i, test := range tests {
		for _, codec := range codecs {
			name := fmt.Sprintf("iteration %d codec %s", i, codec.Name())
			t.Run(name, func(t *testing.T) {
				buf, err := codec.Encode(test)
				require.Nil(t, err)

				str, err := codec.Decode(buf)
				require.Nil(t, err)

				require.Equal(t, test, str)
			})
		}
	}
}
