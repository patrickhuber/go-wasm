package instance

import (
	"github.com/patrickhuber/go-wasm/types"
	"github.com/patrickhuber/go-wasm/values"
)

type Table struct {
	Type    types.Table
	Element []values.Reference
}

func (*Table) instance() {}
