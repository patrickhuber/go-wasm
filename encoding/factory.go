package encoding

import (
	"fmt"
	"os"

	"golang.org/x/text/encoding/unicode"
)

type Factory interface {
	Lookup(Encoding) (Codec, bool)
	Get(Encoding) (Codec, error)
}

type factory struct {
	data map[Encoding]Codec
}

func DefaultFactory() Factory {
	return NewFactory(
		NewUTF16(),
		NewUTF8(),
		NewLatin1(),
		NewUTF16WithEndianess(unicode.BigEndian),
		NewUTF16WithEndianess(unicode.LittleEndian),
	)
}

func NewFactory(codecs ...Codec) Factory {
	codecMap := map[Encoding]Codec{}
	for _, c := range codecs {
		codecMap[c.Encoding()] = c
	}
	return &factory{
		data: codecMap,
	}
}

var ErrNotExist = os.ErrNotExist

func (f *factory) Get(enc Encoding) (Codec, error) {
	codec, ok := f.Lookup(enc)
	if ok {
		return codec, nil
	}
	return nil, fmt.Errorf("encoding %s %w", string(enc), ErrNotExist)
}

func (f *factory) Lookup(enc Encoding) (Codec, bool) {
	c, ok := f.data[enc]
	return c, ok
}
