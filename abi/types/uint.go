package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type U8 struct{}

func (U8) Kind() kind.Kind {
	return kind.U8
}

func (U8) Size() uint32 {
	return 1
}

func (U8) Alignment() uint32 {
	return 1
}

type U16 struct{}

func (U16) Kind() kind.Kind {
	return kind.U16
}

func (U16) Size() uint32 {
	return 2
}

func (U16) Alignment() uint32 {
	return 2
}

type U32 struct{}

func (U32) Kind() kind.Kind {
	return kind.U32
}

func (U32) Size() uint32 {
	return 4
}

func (U32) Alignment() uint32 {
	return 4
}

type U64 struct{}

func (U64) Kind() kind.Kind {
	return kind.U64
}

func (U64) Size() uint32 {
	return 8
}

func (U64) Alignment() uint32 {
	return 8
}
