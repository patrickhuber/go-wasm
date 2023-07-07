package instructions

import "github.com/patrickhuber/go-wasm/types"

type Drop struct{}

func (*Drop) instruction() {}

type Select struct {
	Types []types.Value
}

func (*Select) instruction() {}
