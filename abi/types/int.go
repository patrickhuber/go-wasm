package types

import "github.com/patrickhuber/go-wasm/abi/kind"

const (
	SizeOfS8  = SizeOfU8
	SizeOfS16 = SizeOfU16
	SizeOfS32 = SizeOfU32
	SizeOfS64 = SizeOfU64
)

type S8 struct{}

func (S8) Kind() kind.Kind {
	return kind.S8
}

func (S8) Size() (uint32, error) {
	return SizeOfS8, nil
}

func (S8) Alignment() (uint32, error) {
	return 1, nil
}

func (i S8) Despecialize() ValType {
	return i
}

func (s S8) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type S16 struct{}

func (S16) Kind() kind.Kind {
	return kind.S16
}

func (S16) Size() (uint32, error) {
	return SizeOfS16, nil
}

func (S16) Alignment() (uint32, error) {
	return SizeOfS16, nil
}

func (i S16) Despecialize() ValType {
	return i
}

func (s S16) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type S32 struct{}

func (S32) Kind() kind.Kind {
	return kind.S32
}

func (S32) Size() (uint32, error) {
	return SizeOfS32, nil
}

func (S32) Alignment() (uint32, error) {
	return SizeOfS32, nil
}

func (i S32) Despecialize() ValType {
	return i
}

func (s S32) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}

type S64 struct{}

func (S64) Kind() kind.Kind {
	return kind.S64
}

func (S64) Size() (uint32, error) {
	return SizeOfS64, nil
}

func (S64) Alignment() (uint32, error) {
	return SizeOfS64, nil
}

func (i S64) Despecialize() ValType {
	return i
}

func (s S64) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U64}, nil
}
