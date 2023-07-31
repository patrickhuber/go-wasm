package wasm

type Module struct {
	Sections []Section
}

func (Module) root() {}
