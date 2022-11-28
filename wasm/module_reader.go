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
	module := &Module{}
	sections, err := r.readSections()
	if err != nil {
		return nil, err
	}
	for _, section := range sections {
		switch {
		case section.Code != nil:
			module.Codes = append(module.Codes, section)
		case section.Function != nil:
			module.Functions = append(module.Functions, section)
		case section.Type != nil:
			module.Types = append(module.Types, section)
		}
	}

	return module, nil
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
