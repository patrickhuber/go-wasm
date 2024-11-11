package instance

import (
	"github.com/patrickhuber/go-wasm/address"
	"github.com/patrickhuber/go-wasm/api"
)

type Module struct {
	Directive
	Types             []api.FuncType
	FunctionAddresses []address.Function
	TableAddresses    []address.Table
	MemoryAddresses   []address.Memory
	GlobalAddresses   []address.Global
	ElementAddresses  []address.Element
	DataAddressses    []address.Data
	Exports           []Export
}
