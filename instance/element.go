package instance

import (
	"github.com/patrickhuber/go-wasm/types"
	"github.com/patrickhuber/go-wasm/values"
)

type Element struct {
	Type     types.Reference
	Elements []values.Reference
}

func (*Element) instance() {}
