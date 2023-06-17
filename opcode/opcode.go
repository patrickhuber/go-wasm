package opcode

// Opcode represents the wasm instructions laid out as uint32 using wasm binary encoding. This layout ensures easy execution during dispatch
// https://webassembly.github.io/spec/core/appendix/index-instructions.html#index-instr
type Opcode byte

const (
	Unreachable Opcode = 0x00
	Nop         Opcode = 0x01
	Bloc        Opcode = 0x02

	Return Opcode = 0x0f
	Call   Opcode = 0x10

	Drop Opcode = 0x1a

	LocalGet Opcode = 0x20
	LocalSet Opcode = 0x21

	I32Load Opcode = 0x28
	I64Load Opcode = 0x29
	F32Load Opcode = 0x2A
	F64Load Opcode = 0x2B

	I32Store Opcode = 0x36

	I32Const Opcode = 0x41

	I32Add Opcode = 0x6a
)
