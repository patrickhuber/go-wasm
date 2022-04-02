package wasm

import (
	"io"

	"encoding/binary"
)

type Writer interface {
	Write(module *Module) error
	WriteLebU128(value uint32) error
}

func NewWriter(w io.Writer) Writer {
	return &writer{
		writer: w,
	}
}

type writer struct {
	writer io.Writer
}

func (w *writer) Write(m *Module) error {
	err := binary.Write(w.writer, binary.BigEndian, m.Magic)
	if err != nil {
		return err
	}
	return binary.Write(w.writer, binary.LittleEndian, m.Version)
}

func (w *writer) writeModule(module *Module) error {
	err := binary.Write(w.writer, binary.BigEndian, module.Magic)
	if err != nil {
		return err
	}
	return binary.Write(w.writer, binary.LittleEndian, module.Version)
}

func (w *writer) WriteLebU128(value uint32) error {
	b := make([]byte, 1)
	for {
		b[0] = byte(value & 0b_0111_1111)
		value >>= 7
		if value != 0 {
			b[0] |= 0b_1000_0000
		}
		_, err := w.writer.Write(b)
		if err != nil {
			return err
		}
		if value == 0 {
			break
		}
	}
	return nil
}
