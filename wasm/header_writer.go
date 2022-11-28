package wasm

import (
	"bufio"
	"encoding/binary"
)

type HeaderWriter interface {
	Write(header *Header) error
}

type headerWriter struct {
	writer *bufio.Writer
}

func NewHeaderWriter(writer *bufio.Writer) HeaderWriter {
	return &headerWriter{
		writer: writer,
	}
}

func (w *headerWriter) Write(header *Header) error {
	err := binary.Write(w.writer, binary.BigEndian, header.Magic)
	if err != nil {
		return err
	}
	err = binary.Write(w.writer, binary.LittleEndian, header.Version)
	if err != nil {
		return err
	}
	return binary.Write(w.writer, binary.LittleEndian, header.Layer)
}
