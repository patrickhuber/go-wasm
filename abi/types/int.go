package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type S8 struct{}

func (S8) Kind() kind.Kind {
	return kind.S8
}

func (S8) Size() uint32 {
	return 1
}

func (S8) Alignment() uint32 {
	return 1
}

type S16 struct{}

func (S16) Kind() kind.Kind {
	return kind.S16
}

func (S16) Size() uint32 {
	return 2
}

func (S16) Alignment() uint32 {
	return 2
}

type S32 struct{}

func (S32) Kind() kind.Kind {
	return kind.S32
}

func (S32) Size() uint32 {
	return 4
}

func (S32) Alignment() uint32 {
	return 4
}

type S64 struct{}

func (S64) Kind() kind.Kind {
	return kind.S64
}

func (S64) Size() uint32 {
	return 8
}

func (S64) Alignment() uint32 {
	return 8
}
