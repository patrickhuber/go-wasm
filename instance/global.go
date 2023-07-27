package instance

import (
	"github.com/patrickhuber/go-wasm/types"
	"github.com/patrickhuber/go-wasm/values"
)

type Global struct {
	Type  types.Global
	Value values.Value
}

func (*Global) instance() {}
