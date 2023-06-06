package encoding

import (
	"golang.org/x/text/encoding/unicode"
)

func UTF16() Codec {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return &codec{
		enc:  enc,
		name: "utf-16",
	}
}
