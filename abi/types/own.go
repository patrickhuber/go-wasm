package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Own struct {
	ResourceType *ResourceType
}

func (o *Own) Alignment() (uint32, error) {
	return 4, nil
}

func (*Own) Kind() kind.Kind {
	return kind.Own
}

func (*Own) Size() (uint32, error) {
	return 4, nil
}

func (o *Own) Despecialize() ValType {
	return o
}

func (Own) Flatten() ([]kind.Kind, error) {
	return []kind.Kind{kind.U32}, nil
}
