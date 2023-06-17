package wasm

// Object is the top level element in the wasm ast.
// It can be either a Component or a Module
type Object struct {
	Header    *Header
	Component *Component
	Module    *Module
}
