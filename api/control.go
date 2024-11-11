package api

type BlockType interface {
	blocktype()
}

type BlockTypeIndex struct {
	Index Index
}

func (*BlockTypeIndex) blocktype() {}

type BlockTypeValue struct {
	ValueType ValType
}

func (*BlockTypeValue) blocktype() {}

type Nop struct{}

func (*Nop) instruction() {}

type Unreachable struct{}

func (*Unreachable) instruction() {}

type Block struct {
	Type         BlockType
	Instructions []Instruction
}

func (*Block) instruction() {}

type Loop struct {
	Type         BlockType
	Instructions []Instruction
}

func (*Loop) instruction() {}

type If struct {
	Type         BlockType
	Instructions []Instruction
	Else         *Else
}

func (*If) instruction() {}

type Else struct {
	Instructions Instruction
}

type Branch struct {
	Index LabelIndex
}

func (*Branch) instruction() {}

type BranchIf struct {
	Index LabelIndex
}

func (*BranchIf) instruction() {}

type BranchTable struct {
	Indicies []LabelIndex
	Index    LabelIndex
}

func (*BranchTable) instruction() {}

type Return struct{}

func (*Return) instruction() {}

type Call struct {
	Index FunctionIndex
}

func (*Call) instruction() {}

type CallIndirect struct {
	Table TableIndex
	Type  TypeIndex
}

func (*CallIndirect) instruction() {}

type End struct {
	Instruction
}
