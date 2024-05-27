package runtime_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/go-wasm/engine"
	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/runtime"
	"github.com/patrickhuber/go-wasm/store"
)

func TestInstantiate(t *testing.T) {

	e := engine.New()
	s := store.New(e)
	r := runtime.New(s)

	add, err := os.Open("../fixtures/add/add.wasm")
	if err != nil {
		t.Fatal(err)
	}

	d, err := r.Instantiate(add)
	if err != nil {
		t.Fatal(err)
	}

	m, ok := d.(*instance.Module)
	if !ok {
		t.Fatal("expected module instance")
	}

	if m == nil {
		t.Fatalf("module instance is nil")
	}
	if len(m.FunctionAddresses) == 0 {
		t.Fatalf("expected non empty function address list")
	}
}
