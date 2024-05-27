package instruction

import "github.com/patrickhuber/go-wasm/indicies"

type LocalGet struct {
	Index indicies.Local
}

func (LocalGet) instruction() {}

type LocalSet struct {
	Index indicies.Local
}

func (LocalSet) instruction() {}

type LocalTee struct {
	Index indicies.Local
}

func (LocalTee) instruction() {}

type GlobalGet struct {
	Index indicies.Global
}

func (GlobalGet) instruction() {}

type GlobalSet struct {
	Index indicies.Global
}

func (GlobalSet) instruction() {}
