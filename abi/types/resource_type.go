package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type DTorFunc func(int)

type ResourceType struct {
	Impl *ComponentInstance
	DTor *DTorFunc
}

func (rt *ResourceType) Kind() kind.Kind {
	return kind.ResourceType
}
