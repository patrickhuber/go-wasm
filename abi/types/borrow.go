package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Borrow struct {
	ResourceType *ResourceType
}

func (*Borrow) Alignment() uint32 {
	return 4
}

func (*Borrow) Kind() kind.Kind {
	return kind.Borrow
}

func (*Borrow) Size() uint32 {
	return 4
}

func (b *Borrow) Despecialize() ValType {
	return b
}
