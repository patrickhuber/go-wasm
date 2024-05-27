package instruction

import "github.com/patrickhuber/go-wasm/indicies"

type MemoryArg struct {
	Offset uint32
	Align  uint32
}

type Int32Load struct {
	MemoryArg
}

func (*Int32Load) instruction() {}

type Int32Store struct {
	MemoryArg
}

func (*Int32Store) instruction() {}

type F32Load struct {
	MemoryArg
}

func (*F32Load) instruction() {}

type I32Load8 struct {
	MemoryArg
}

func (*I32Load8) instruction() {}

type U32Load8u struct {
	MemoryArg
}

func (*U32Load8u) instruction() {}

type I32Load16 struct {
	MemoryArg
}

func (*I32Load16) instruction() {}

type U32Load16 struct {
	MemoryArg
}

func (*U32Load16) instruction() {}

type I64Load32 struct {
	MemoryArg
}

func (*I64Load32) instruction() {}

type U64Load32 struct {
	MemoryArg
}

func (*U64Load32) instruction() {}

type I32Store8 struct {
	MemoryArg
}

func (*I32Store8) instruction() {}

type U32Store8 struct {
	MemoryArg
}

func (*U32Store8) instruction() {}

type I32Store16 struct {
	MemoryArg
}

func (*I32Store16) instruction() {}

type U32Store16 struct {
	MemoryArg
}

func (*U32Store16) instruction() {}

type I64Store32 struct {
	MemoryArg
}

func (*I64Store32) instruction() {}

type U64Store32 struct {
	MemoryArg
}

func (*U64Store32) instruction() {}

type MemorySize struct{}

func (*MemorySize) instruction() {}

type MemoryGrow struct{}

func (*MemoryGrow) instruction() {}

type MemoryCopy struct{}

func (*MemoryCopy) instruction() {}

type MemoryInit struct {
	Index indicies.Data
}

func (*MemoryInit) instruction() {}

type DataDrop struct {
	Index indicies.Data
}

func (*DataDrop) instruction() {}
