package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type StringEncoding string

const (
	None                StringEncoding = "none"
	Utf8                StringEncoding = "utf8"
	Utf16               StringEncoding = "utf16"
	Latin1Utf16         StringEncoding = "latin1+utf16"
	MaxStringByteLength uint32         = (1 << 31) - 1
)

type String struct{}

func (String) Kind() kind.Kind {
	return kind.String
}

func (String) Size() uint32 {
	return 8
}

func (String) Alignment() uint32 {
	return 4
}

func (s String) Despecialize() ValType {
	return s
}

func (String) Flatten() []kind.Kind {
	return []kind.Kind{kind.S32, kind.S32}
}
