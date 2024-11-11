package instance

import "github.com/patrickhuber/go-wasm/api"

type Function interface {
}

type ModuleFunction struct {
	Type api.FuncType
}

func (*ModuleFunction) instance() {}

type HostCodeFunction struct {
	Type     api.FuncType
	HostCode *HostFunction
}

func (*HostCodeFunction) instance() {}

type HostFunction struct{}
