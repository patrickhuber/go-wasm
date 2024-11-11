package instance

import "github.com/patrickhuber/go-wasm/api"

type Memory struct {
	Type api.Mem
	Data []byte
}

func (*Memory) instance() {}
