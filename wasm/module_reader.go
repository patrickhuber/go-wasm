package wasm

import (
	"bufio"
	"io"
)

type ModuleReader interface {
	Read() (*Module, error)
}

func NewModuleReader(r io.Reader) ModuleReader {
	return &moduleReader{
		reader: bufio.NewReader(r),
	}
}

type moduleReader struct {
	reader *bufio.Reader
}

func (r *moduleReader) Read() (*Module, error) {

	sections, err := r.readSections()
	if err != nil {
		return nil, err
	}
	return &Module{
		Sections: sections,
	}, nil
}

func (r *moduleReader) readSections() ([]Section, error) {
	var sections []Section
	reader := NewSectionReader(r.reader, LayerCore)
	for {
		section, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return sections, nil
			}
			return nil, err
		}
		sections = append(sections, *section)
	}
}
