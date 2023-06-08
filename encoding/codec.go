package encoding

import (
	"io"

	"golang.org/x/text/encoding"
)

type Codec interface {
	Encoder
	Decoder
	Encoding() Encoding
	Alignment() int
	RuneSize() int
}

type Encoder interface {
	Encode(dst io.Writer, src io.Reader) error
}

type Decoder interface {
	Decode(dst io.Writer, src io.Reader) error
}

type Encoding string

const (
	None Encoding = ""
)

type codec struct {
	name      Encoding
	enc       encoding.Encoding
	alignment int
	runeSize  int
}

func (c *codec) Encode(dst io.Writer, src io.Reader) error {
	encoder := c.enc.NewEncoder()
	writer := encoder.Writer(dst)
	_, err := io.Copy(writer, src)
	return err
}

func (c *codec) Decode(dst io.Writer, src io.Reader) error {
	decoder := c.enc.NewDecoder()
	reader := decoder.Reader(src)
	_, err := io.Copy(dst, reader)
	return err
}

func (c *codec) Encoding() Encoding {
	return c.name
}

func (c *codec) Alignment() int {
	return c.alignment
}

func (c *codec) RuneSize() int {
	return c.runeSize
}
