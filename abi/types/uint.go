package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type U8 struct{}

const (
	SizeOfU8  = 1
	SizeOfU16 = 2
	SizeOfU32 = 4
	SizeOfU64 = 8
)

func (U8) Kind() kind.Kind {
	return kind.U8
}

func (U8) Size() (uint32, error) {
	return SizeOfU8, nil
}

func (U8) Alignment() (uint32, error) {
	return 1, nil
}

func (u U8) Despecialize() ValType {
	return u
}

func (u U8) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type U16 struct{}

func (U16) Kind() kind.Kind {
	return kind.U16
}

func (U16) Size() (uint32, error) {
	return SizeOfU16, nil
}

func (U16) Alignment() (uint32, error) {
	return 2, nil
}

func (u U16) Despecialize() ValType {
	return u
}

func (u U16) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type U32 struct{}

func (U32) Kind() kind.Kind {
	return kind.U32
}

func (U32) Size() (uint32, error) {
	return SizeOfU32, nil
}

func (U32) Alignment() (uint32, error) {
	return 4, nil
}

func (u U32) Despecialize() ValType {
	return u
}

func (u U32) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type U64 struct{}

func (U64) Kind() kind.Kind {
	return kind.U64
}

func (U64) Size() (uint32, error) {
	return SizeOfU64, nil
}

func (U64) Alignment() (uint32, error) {
	return 8, nil
}

func (u U64) Despecialize() ValType {
	return u
}

func (u U64) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U64}, nil
}
