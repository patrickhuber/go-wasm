package encoding

import "golang.org/x/text/encoding/unicode"

func UTF8() Codec {
	enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	return &codec{
		enc:  enc,
		name: "utf-8",
	}
}

type utf8Codec struct {
}

func (c *utf8Codec) Encode(str string) ([]byte, error) {
	return []byte(str), nil
}

func (c *utf8Codec) Decode(b []byte) (string, error) {
	return string(b), nil
}

func (c *utf8Codec) Name() string {
	return "utf-8"
}
