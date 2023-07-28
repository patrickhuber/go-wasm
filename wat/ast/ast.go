package ast

type Ast interface {
	ast()
}

type Component struct {
}

func (*Component) ast() {}

type Module struct {
	Functions []Section
	Memory    []Section
}

func (*Module) ast() {}

type Section struct {
	Function *Function
	Memory   *Memory
}

type Function struct {
	ID           *Identifier
	Parameters   []Parameter
	Results      []Result
	Instructions []Instruction
}

type Memory struct {
	ID     *Identifier
	Limits Limits
}

type Limits struct {
	HasMax bool
	Min    uint32
	Max    uint32
}

type Identifier string

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

type Type string

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
