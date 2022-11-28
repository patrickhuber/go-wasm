package wasm

type Instruction struct {
	OpCode OpCode
	Local  *LocalInstruction
	Const  *ConstInstruction
}

//go:generate stringer -type=OpCode

// OpCode represents an instruction within wasm
type OpCode byte

const (
	Unreachable OpCode = 0x00
	Nop         OpCode = 0x01
	End         OpCode = 0x0B
	LocalGet    OpCode = 0x20
	LocalSet    OpCode = 0x21
	LocalTee    OpCode = 0x22
	GlobalGet   OpCode = 0x23
	GlobalSet   OpCode = 0x24

	I32Const OpCode = 0x41
	I64Const OpCode = 0x42
	F32Const OpCode = 0x43
	F64Const OpCode = 0x44

	I32Clz    OpCode = 0x67
	I32Ctz    OpCode = 0x68
	I32PopCnt OpCode = 0x69
	I32Add    OpCode = 0x6A
)

type LocalInstruction struct {
	Index *uint32
	Tag   *string
}

type LocalGlobal struct {
	Index *uint32
	Tag   *string
}

type ConstInstruction struct {
	I32 *int32
	I64 *int64
	F32 *float32
	F64 *float64
}
