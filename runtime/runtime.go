package runtime

import (
	"io"

	"github.com/patrickhuber/go-wasm/instance"
	"github.com/patrickhuber/go-wasm/store"
)

type Runtime interface {
	Instantiate(wasm io.Reader) (instance.Directive, error)
}

type runtime struct {
	store *store.Store
}

func New(store *store.Store) Runtime {
	return &runtime{
		store: store,
	}
}

func (r *runtime) Instantiate(wasm io.Reader) (instance.Directive, error) {

	// document, err := binary.Read(wasm)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}
