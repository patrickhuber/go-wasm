package runtime_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/runtime"
)

func TestMachine(t *testing.T) {
	machine := runtime.NewMachine()
	module := &api.Module{}
	externals := []api.External{}
	err := machine.Execute(module, externals)
	if err != nil {
		t.Error(err)
	}
}
