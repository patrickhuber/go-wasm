package wasm

import (
	"bufio"
	"io"
)

type ModuleWriter interface {
	Write(module *Module) error
}

func NewModuleWriter(w io.Writer) ModuleWriter {
	return &moduleWriter{
		writer: bufio.NewWriter(w),
	}
}

type moduleWriter struct {
	writer *bufio.Writer
}

func (w *moduleWriter) Write(m *Module) error {

	return nil
}
