package encoding

import "golang.org/x/text/encoding/unicode"

const (
	UTF8 Encoding = "utf-8"
)

func NewUTF8() Codec {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return &codec{
		enc:       enc,
		name:      UTF8,
		alignment: 1,
		runeSize:  1,
	}
}
