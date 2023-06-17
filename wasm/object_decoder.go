package wasm

import (
	"bufio"
	"fmt"
	"io"
)

type ObjectDecoder interface {
	Decode() (*Object, error)
}

type objectReader struct {
	reader *bufio.Reader
}

func NewObjectDecoder(reader io.Reader) ObjectDecoder {
	return &objectReader{
		reader: bufio.NewReader(reader),
	}
}

func (r *objectReader) Decode() (*Object, error) {
	header, err := NewHeaderReader(r.reader).Read()
	if err != nil {
		return nil, err
	}
	object := &Object{
		Header: header,
	}
	switch {
	case header.Layer == LayerComponent && header.Version == VersionExperimental:
		component, err := NewComponentReader(r.reader).Read()
		if err != nil {
			return nil, err
		}
		object.Component = component

	case header.Layer == LayerCore && header.Version == Version1:
		module, err := NewModuleReader(r.reader).Read()
		if err != nil {
			return nil, err
		}
		object.Module = module

	default:
		return nil, fmt.Errorf("invalid version")
	}
	return object, nil
}
