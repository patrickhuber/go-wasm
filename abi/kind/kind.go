package kind

type Kind int

//go:generate stringer -type=Kind
const (
	U32 Kind = iota
	U64
	Float32
	Float64
)
