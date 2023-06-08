package encoding

import (
	"golang.org/x/text/encoding/unicode"
)

const (
	UTF16   Encoding = "utf-16"
	UTF16BE Encoding = "utf-16-be"
	UTF16LE Encoding = "utf-16-le"
)

func NewUTF16() Codec {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return &codec{
		enc:       enc,
		name:      UTF16,
		alignment: 2,
		runeSize:  2,
	}
}

func NewUTF16WithEndianess(endianess unicode.Endianness) Codec {
	enc := unicode.UTF16(endianess, unicode.IgnoreBOM)
	name := UTF16BE
	if endianess == unicode.LittleEndian {
		name = UTF16LE
	}
	return &codec{
		enc:       enc,
		name:      name,
		alignment: 2,
		runeSize:  2,
	}
}
