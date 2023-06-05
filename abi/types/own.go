package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Own struct {
	ResourceType *ResourceType
}

func (o *Own) Alignment() uint32 {
	return 4
}

func (*Own) Kind() kind.Kind {
	return kind.Own
}

func (*Own) Size() uint32 {
	return 4
}

func (o *Own) Despecialize() ValType {
	return o
}

func (Own) Flatten() []kind.Kind {
	return []kind.Kind{kind.S32}
}
