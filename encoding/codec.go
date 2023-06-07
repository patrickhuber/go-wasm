package encoding

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/text/encoding"
)

type Codec interface {
	Encoder
	Decoder
	Encoding() Encoding
}

type Encoder interface {
	Encode(src string) ([]byte, error)
}

type Decoder interface {
	Decode(src []byte) (string, error)
}

type Encoding string

type codec struct {
	name Encoding
	enc  encoding.Encoding
}

func (c *codec) Encode(src string) ([]byte, error) {
	encoder := c.enc.NewEncoder()
	reader := strings.NewReader(src)
	buf := &bytes.Buffer{}
	writer := encoder.Writer(buf)
	_, err := io.Copy(writer, reader)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *codec) Decode(src []byte) (string, error) {
	decoder := c.enc.NewDecoder()
	reader := decoder.Reader(bytes.NewReader(src))
	writer := &strings.Builder{}
	_, err := io.Copy(writer, reader)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}

func (c *codec) Encoding() Encoding {
	return c.name
}
