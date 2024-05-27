package instruction

import (
	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/types"
)

type BlockType interface {
	blocktype()
}

type BlockTypeIndex struct {
	Index indicies.Index
}

func (*BlockTypeIndex) blocktype() {}

type BlockTypeValue struct {
	ValueType types.Value
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
	Index indicies.Label
}

func (*Branch) instruction() {}

type BranchIf struct {
	Index indicies.Label
}

func (*BranchIf) instruction() {}

type BranchTable struct {
	Indicies []indicies.Label
	Index    indicies.Label
}

func (*BranchTable) instruction() {}

type Return struct{}

func (*Return) instruction() {}

type Call struct {
	Index indicies.Function
}

func (*Call) instruction() {}

type CallIndirect struct {
	Table indicies.Table
	Type  indicies.Type
}

func (*CallIndirect) instruction() {}

type End struct {
	Instruction
}
