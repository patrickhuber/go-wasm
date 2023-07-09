package encoding

import (
	"fmt"

	match "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"golang.org/x/text/encoding/unicode"
)

type Factory interface {
	Lookup(Encoding) match.Option[Codec]
	Get(Encoding) match.Result[Codec]
}

type factory struct {
	data map[Encoding]Codec
}

func DefaultFactory() Factory {
	return NewFactory(
		NewUTF16(),
		NewUTF8(),
		NewLatin1(),
		NewLatin1Utf16(),
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

var ErrNotExist = fmt.Errorf("does not exist")

func (f *factory) Get(enc Encoding) (res match.Result[Codec]) {
	switch codec := f.Lookup(enc).(type) {
	case match.Some[Codec]:
		return result.Ok(codec.Value())
	default:
		return result.Errorf[Codec]("encoding %s %w", string(enc), ErrNotExist)
	}
}

func (f *factory) Lookup(enc Encoding) match.Option[Codec] {
	c, ok := f.data[enc]
	return option.New(c, ok)
}
