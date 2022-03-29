package wasm

type Module struct {
	Functions []Function
}

type Section struct {
	Function *Function
}

type Function struct {
	ID           *Identifier
	Parameters   []Parameter
	Results      []Result
	Instructions []Instruction
}

type Identifier string

func Pointer[T any](item T) *T {
	return &item
}

type Parameter struct {
	ID   *Identifier
	Type Type
}

type Local struct {
	Type Type
}

type Result struct {
	Type Type
}

type Type int

const (
	I32 Type = iota
	I64
)

type Instruction struct {
	Block *Block
	Plain *Plain
}

type Block struct{}
type Plain struct {
	Local *LocalInstruction
	I32   *I32Instruction
}

type LocalOperation string

const (
	LocalGet LocalOperation = "get"
	LocalSet LocalOperation = "set"
	LocalTee LocalOperation = "tee"
)

type LocalInstruction struct {
	Operation LocalOperation
	ID        *Identifier
	Index     *int
}

type I32Instruction struct {
	Operation BinaryOperation
}

type BinaryOperation string

const (
	BinaryOperationAdd BinaryOperation = "add"
	BinaryOperationSub BinaryOperation = "sub"
	BinaryOperationMul BinaryOperation = "mul"
)
