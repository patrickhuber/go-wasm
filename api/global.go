package api

type Mutable int

const (
	Const Mutable = iota
	Var
)

type Global struct {
	Mutable Mutable
	Value   ValType
}

// external implements External interface
func (*Global) external() {}
