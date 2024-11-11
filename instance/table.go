package instance

import (
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/values"
)

type Table struct {
	Type    api.Table
	Element []values.Reference
}

func (*Table) instance() {}
