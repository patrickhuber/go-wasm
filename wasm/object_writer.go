package wasm

import (
	"bufio"
	"io"
)

type ObjectWriter interface {
	Write(object *Object) error
}

type objectWriter struct {
	writer *bufio.Writer
}

func NewObjectWriter(writer io.Writer) ObjectWriter {
	return &objectWriter{
		writer: bufio.NewWriter(writer),
	}
}

func (w *objectWriter) Write(object *Object) error {
	err := NewHeaderWriter(w.writer).Write(object.Header)
	if err != nil {
		return err
	}
	return w.writer.Flush()
}
