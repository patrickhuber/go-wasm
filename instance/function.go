package instance

import "github.com/patrickhuber/go-wasm/types"

type Function interface {
	Instance
}

type ModuleFunction struct {
	Type types.Function
}

func (*ModuleFunction) instance() {}

type HostCodeFunction struct {
	Type     types.Function
	HostCode *HostFunction
}

func (*HostCodeFunction) instance() {}

type HostFunction struct{}
