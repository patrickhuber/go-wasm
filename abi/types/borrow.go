package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Borrow struct {
	ResourceType *ResourceType
}

func (Borrow) Alignment() (uint32, error) {
	return 4, nil
}

func (Borrow) Kind() kind.Kind {
	return kind.Borrow
}

func (Borrow) Size() (uint32, error) {
	return 4, nil
}

func (b *Borrow) Despecialize() ValType {
	return b
}

func (Borrow) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}
