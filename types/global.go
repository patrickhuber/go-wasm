package types

type Mutable int

const (
	Const Mutable = iota
	Var
)

type Global struct {
	Mutable Mutable
	Value   Value
}

// external implements External interface
func (*Global) external() {}
