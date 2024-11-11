package instance

import (
	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/values"
)

type Element struct {
	Type     api.Reference
	Elements []values.Reference
}

func (*Element) instance() {}
