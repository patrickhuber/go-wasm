package kind

type Kind int

//go:generate stringer -type=Kind
const (
	Bool Kind = iota
	U8
	S8
	U16
	S16
	U32
	S32
	U64
	S64
	Float32
	Float64
	Char
	String
	List
	Record
	Tuple
	Variant
	Enum
	Union
	Option
	Result
	Flags
	Own
	Borrow
	ValType
	ExternType
	CoreExternType
	ResourceType
)
