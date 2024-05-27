package instance

import (
	"github.com/patrickhuber/go-wasm/address"
	"github.com/patrickhuber/go-wasm/types"
)

type Module struct {
	Directive
	Types             []types.Function
	FunctionAddresses []address.Function
	TableAddresses    []address.Table
	MemoryAddresses   []address.Memory
	GlobalAddresses   []address.Global
	ElementAddresses  []address.Element
	DataAddressses    []address.Data
	Exports           []Export
}
