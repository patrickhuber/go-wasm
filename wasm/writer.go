package wasm

import (
	"bufio"
	"io"

	"encoding/binary"

	"github.com/patrickhuber/go-wasm/leb128"
)

type Writer interface {
	Write(module *Module) error
}

func NewWriter(w io.Writer) Writer {
	return &writer{
		writer: bufio.NewWriter(w),
	}
}

type writer struct {
	writer *bufio.Writer
}

func (w *writer) Write(m *Module) error {
	err := binary.Write(w.writer, binary.BigEndian, m.Magic)
	if err != nil {
		return err
	}
	err = binary.Write(w.writer, binary.LittleEndian, m.Version)
	if err != nil {
		return err
	}

	return w.writer.Flush()
}

func (w *writer) writeModule(module *Module) error {
	err := binary.Write(w.writer, binary.BigEndian, module.Magic)
	if err != nil {
		return err
	}
	return binary.Write(w.writer, binary.LittleEndian, module.Version)
}

func (w *writer) writeLebU128(value uint32) error {
	_, err := leb128.Encode(w.writer, value)
	return err
}
