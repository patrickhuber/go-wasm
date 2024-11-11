package api

type LocalGet struct {
	Index LocalIndex
}

func (LocalGet) instruction() {}

type LocalSet struct {
	Index LocalIndex
}

func (LocalSet) instruction() {}

type LocalTee struct {
	Index LocalIndex
}

func (LocalTee) instruction() {}

type GlobalGet struct {
	Index GlobalIndex
}

func (GlobalGet) instruction() {}

type GlobalSet struct {
	Index GlobalIndex
}

func (GlobalSet) instruction() {}
