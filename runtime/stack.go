package runtime

import "github.com/patrickhuber/go-wasm/wasm"

type Stack interface {
	Push(StackItem)
	Pop() StackItem
	I32Stack
	I64Stack
	F32Stack
	F64Stack
}

type I32Stack interface {
	PopUint32() uint32
	PushUint32(uint32)
}

type I64Stack interface {
	PopUint64() uint64
	PushUint64(uint64)
}

type F32Stack interface {
	PopFloat32() float32
	PushFloat32(float32)
}

type F64Stack interface {
	PopFloat64() float64
	PushFloat64(float64)
}

type StackItem struct {
	Value      *Value
	Label      *Label
	Activation *Activation
}

type Number struct {
	OpCode wasm.OpCode
	Value  []byte
}

type Vec struct {
}

type Ref struct {
}

type Label struct{}

type Activation struct {
	Frames []Frame
}

type Frame struct {
	Locals []Value
	Module *ModuleInstance
}
