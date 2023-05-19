package runtime

import (
	"encoding/binary"

	gstack "github.com/patrickhuber/go-collections/generic/stack"
	"github.com/patrickhuber/go-wasm/wasm"
)

type stack struct {
	inner gstack.Stack[StackItem]
}

func NewStack() Stack {
	return &stack{
		inner: gstack.New[StackItem](),
	}
}

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
	Value      *wasm.Value
	Label      *Label
	Activation *Activation
}

type Label struct{}

type Activation struct {
	Frames []Frame
}

type Frame struct {
	Locals []wasm.Value
	Module *ModuleInstance
}

func (s *stack) Push(item StackItem) {
	s.inner.Push(item)
}

func (s *stack) Pop() StackItem {
	return s.inner.Pop()
}

// PopUint32 implements Stack
func (s *stack) PopUint32() uint32 {
	item := s.Pop()
	return item.Value.Uint32()
}

// PushUint32 implements Stack
func (s *stack) PushUint32(value uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	item := StackItem{
		Value: &wasm.Value{
			Number: &wasm.Number{
				OpCode: wasm.I32Const,
				Value:  b,
			},
		},
	}
	s.Push(item)
}

// PopUint64 implements Stack
func (s *stack) PopUint64() uint64 {
	panic("unimplemented")
}

// PushUint64 implements Stack
func (s *stack) PushUint64(uint64) {
	panic("unimplemented")
}

// PopFloat32 implements Stack
func (s *stack) PopFloat32() float32 {
	panic("unimplemented")
}

// PushFloat32 implements Stack
func (s *stack) PushFloat32(float32) {
	panic("unimplemented")
}

// PopFloat64 implements Stack
func (s *stack) PopFloat64() float64 {
	panic("unimplemented")
}

// PushFloat64 implements Stack
func (s *stack) PushFloat64(float64) {
	panic("unimplemented")
}
