package instance

import "github.com/patrickhuber/go-wasm/types"

type Memory struct {
	Type types.Memory
	Data []byte
}

func (*Memory) instance() {}
