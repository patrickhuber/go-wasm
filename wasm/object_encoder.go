package wasm

import (
	"bufio"
	"io"
)

type ObjectEncoder interface {
	Encode(object *Object) error
}

type objectEncoder struct {
	writer *bufio.Writer
}

func NewObjectEncoder(writer io.Writer) ObjectEncoder {
	return &objectEncoder{
		writer: bufio.NewWriter(writer),
	}
}

func (w *objectEncoder) Encode(object *Object) error {
	err := NewHeaderWriter(w.writer).Write(object.Header)
	if err != nil {
		return err
	}
	return w.writer.Flush()
}
