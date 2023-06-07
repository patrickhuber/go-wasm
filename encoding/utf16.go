package encoding

import (
	"golang.org/x/text/encoding/unicode"
)

const (
	UTF16 Encoding = "utf-16"
)

func NewUTF16() Codec {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return &codec{
		enc:  enc,
		name: "utf-16",
	}
}
