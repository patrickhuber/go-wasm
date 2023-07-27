package instance

import "github.com/patrickhuber/go-wasm/address"

type Export struct {
	Name  string
	Value address.ExternalValue
}

func (*Export) instance() {}
