package runtime

import "github.com/patrickhuber/go-wasm/instance"

// Store represents all global state
// see https://webassembly.github.io/spec/core/exec/runtime.html#store
type Store struct {
	Funcs   []instance.Function
	Tables  []instance.Table
	Mems    []instance.Memory
	Globals []instance.Global
	Elems   []instance.Element
	Datas   []instance.Data
	Modules []instance.Module
}
