package kind

type Kind int

//go:generate stringer -type=Kind
const (
	U8 Kind = iota
	U16
	U32
	U64
	Float32
	Float64
)
