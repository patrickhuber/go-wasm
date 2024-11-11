package instance

import (
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/values"
)

type Global struct {
	Type  api.Global
	Value values.Value
}

func (*Global) instance() {}
