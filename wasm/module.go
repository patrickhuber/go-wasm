package wasm

type Module struct {
	Functions []Section
	Types     []Section
	Codes     []Section
}
