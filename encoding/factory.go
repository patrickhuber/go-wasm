package encoding

type Factory interface {
	Get(Encoding) (Codec, bool)
}

type factory struct {
	data map[Encoding]Codec
}

func DefaultFactory() Factory {
	return NewFactory(
		NewUTF16(),
		NewUTF8(),
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

func (f *factory) Get(enc Encoding) (Codec, bool) {
	c, ok := f.data[enc]
	return c, ok
}
