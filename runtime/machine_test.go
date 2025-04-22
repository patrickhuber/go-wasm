package runtime_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/runtime"
	"github.com/patrickhuber/go-wasm/values"
)

func TestMachine(t *testing.T) {
	machine := runtime.NewMachine()
	module := &api.Module{}
	externals := []values.Value{}
	returns, err := machine.Execute(module, externals)
	if err != nil {
		t.Error(err)
	}
	if len(returns) != 0 {
		t.Errorf("expected 0 returns but found %d", len(returns))
	}
}
