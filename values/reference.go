package values

import "github.com/patrickhuber/go-wasm/address"

type Reference interface {
	reference()
	Value
}

type NullReference struct{}

func (*NullReference) reference() {}
func (*NullReference) value()     {}

type AddressReference struct{}

func (*AddressReference) reference() {}
func (*AddressReference) value()     {}

type ExternalReference struct {
	Address address.External
}

func (*ExternalReference) reference() {}
func (*ExternalReference) value()     {}
