package encoding

import "golang.org/x/text/encoding/unicode"

const (
	UTF8 Encoding = "utf-8"
)

func NewUTF8() Codec {
	return &codec{
		enc:       unicode.UTF8,
		name:      UTF8,
		alignment: 1,
		runeSize:  1,
	}
}
