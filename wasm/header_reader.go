package wasm

import (
	"encoding/binary"
	"io"
)

type HeaderReader interface {
	Read() (*Header, error)
}

func NewHeaderReader(reader io.Reader) HeaderReader {
	return &headerReader{
		reader: reader,
	}
}

type headerReader struct {
	reader io.Reader
}

func (r *headerReader) Read() (*Header, error) {
	header := &Header{}
	err := binary.Read(r.reader, binary.BigEndian, &header.Magic)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r.reader, binary.LittleEndian, &header.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r.reader, binary.LittleEndian, &header.Layer)
	return header, err
}
