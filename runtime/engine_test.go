package runtime_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/runtime"
	"github.com/patrickhuber/go-wasm/store"
	"github.com/stretchr/testify/require"
)

func TestEngineRunAdd(t *testing.T) {
	stack := runtime.NewStack()
	store := store.Store{}
	engine := runtime.NewEngine(stack, store)
	require.Nil(t, engine.Instantiate())
	require.Nil(t, engine.Invoke())
}
