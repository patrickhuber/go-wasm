package encoding

import "golang.org/x/text/encoding/charmap"

const (
	Latin1      Encoding = "latin1"
	Latin1Utf16 Encoding = "latin1+utf16"
)

func NewLatin1() Codec {
	return &codec{
		enc:       charmap.ISO8859_1,
		alignment: 2,
		name:      Latin1,
		runeSize:  1,
	}
}
