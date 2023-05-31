package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type StringEncoding string

const (
	Utf8        StringEncoding = "utf8"
	Utf16       StringEncoding = "utf16"
	Latin1Utf16 StringEncoding = "latin1+utf16"
)

type String string

func (*String) Kind() kind.Kind {
	return kind.String
}

func (*String) Size() uint32 {
	return 8
}

func (*String) Alignment() uint32 {
	return 4
}
