package wasm

import (
	"bufio"
	"io"
)

type ComponentReader interface {
	Read() (*Component, error)
}

type componentReader struct {
	reader *bufio.Reader
}

func NewComponentReader(reader *bufio.Reader) ComponentReader {
	return &componentReader{
		reader: reader,
	}
}

func (r *componentReader) Read() (*Component, error) {
	component := &Component{}
	sections, err := r.readSections()
	if err != nil {
		return nil, err
	}
	for _, section := range sections {
		switch {
		case section.Custom != nil:
			component.Custom = append(component.Custom, section)
		}
	}
	return component, nil
}

func (r *componentReader) readSections() ([]Section, error) {
	var sections []Section
	reader := NewSectionReader(r.reader, LayerComponent)
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
