package binary

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"encoding/binary"

	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/instruction"
	"github.com/patrickhuber/go-wasm/leb128"
	"github.com/patrickhuber/go-wasm/opcode"
)

func Read(reader io.Reader) (*Document, error) {

	preamble, err := ReadPreamble(reader)
	if err != nil {
		return nil, err
	}

	var root Root
	switch preamble.Version {
	case ComponentVersion:
		return nil, fmt.Errorf("component binary format not supported yet")
	case ModuleVersion:
		root, err = ReadModule(reader)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid version %d", preamble.Version)
	}

	return &Document{
		Preamble: preamble,
		Root:     root,
	}, nil
}

func ReadPreamble(reader io.Reader) (*Preamble, error) {
	preamble := &Preamble{}

	var err error

	preamble.Magic, err = read[[4]byte](reader)
	if err != nil {
		return nil, err
	}

	if !slices.Equal(Magic, preamble.Magic[0:]) {
		return nil, fmt.Errorf("expected magic %v found %v", Magic, preamble.Magic)
	}

	preamble.Version, err = ReadUInt16(reader)
	if err != nil {
		return nil, err
	}

	preamble.Layer, err = ReadUInt16(reader)
	if err != nil {
		return nil, err
	}

	return preamble, nil
}

func ReadModule(reader io.Reader) (*Module, error) {
	var sections []Section
	for {
		section, err := ReadSection(reader)
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			}
			return nil, err
		}
		sections = append(sections, section)
	}
	return &Module{
		Sections: sections,
	}, nil
}

func ReadSection(reader io.Reader) (Section, error) {

	id, err := ReadByte(reader)
	if err != nil {
		return nil, err
	}

	// read the size
	size, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	switch SectionID(id) {
	case TypeSectionID:
		return ReadTypeSection(size, reader)
	case FunctionSectionID:
		return ReadFunctionSection(size, reader)
	case CodeSectionID:
		return ReadCodeSection(size, reader)
	}
	return nil, fmt.Errorf("invalid section id %d", id)
}

func ReadTypeSection(size uint32, reader io.Reader) (*TypeSection, error) {
	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}
	types := make([]*FunctionType, count)
	for i := uint32(0); i < count; i++ {
		t, err := ReadType(reader)
		if err != nil {
			return nil, err
		}
		types[i] = t
	}
	return &TypeSection{
		ID:    TypeSectionID,
		Size:  size,
		Types: types,
	}, nil
}

func ReadType(reader io.Reader) (*FunctionType, error) {
	b, err := ReadByte(reader)
	if err != nil {
		return nil, err
	}
	if b != 0x60 {
		return nil, fmt.Errorf("expected byte 0x60 but found %b", b)
	}
	parameters, err := ReadValueTypeVector(reader)
	if err != nil {
		return nil, err
	}
	results, err := ReadValueTypeVector(reader)
	if err != nil {
		return nil, err
	}
	return &FunctionType{
		Parameters: parameters,
		Results:    results,
	}, nil
}

func ReadValueTypeVector(reader io.Reader) ([]ValueType, error) {
	// read the size of the vector
	size, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}
	valueTypes := make([]ValueType, size)
	for i := uint32(0); i < size; i++ {
		vt, err := ReadValueType(reader)
		if err != nil {
			return nil, err
		}
		valueTypes[i] = vt
	}
	return valueTypes, nil
}

func ReadValueType(reader io.Reader) (ValueType, error) {
	b, err := ReadByte(reader)
	if err != nil {
		return 0, err
	}
	switch ValueType(b) {
	case I32:
		return I32, nil
	case I64:
		return I64, nil
	case F32:
		return F32, nil
	case F64:
		return F64, nil
	}
	return 0, fmt.Errorf("invalid ValueType found %b", b)
}

func ReadFunctionSection(size uint32, reader io.Reader) (*FunctionSection, error) {

	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	types := make([]uint32, count)
	for i := uint32(0); i < count; i++ {
		index, err := ReadLebU128(reader)
		if err != nil {
			return nil, err
		}
		types[i] = index
	}

	return &FunctionSection{
		ID:    FunctionSectionID,
		Size:  size,
		Types: types,
	}, nil
}

func ReadCodeSection(size uint32, reader io.Reader) (*CodeSection, error) {
	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	codes := make([]*Code, count)
	for i := uint32(0); i < count; i++ {
		code, err := ReadCode(reader)
		if err != nil {
			return nil, err
		}
		codes[i] = code
	}
	return &CodeSection{
		ID:    CodeSectionID,
		Size:  size,
		Codes: codes,
	}, nil
}

func ReadCode(reader io.Reader) (*Code, error) {

	size, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	locals := make([]Local, count)
	for i := uint32(0); i < count; i++ {
		local, err := ReadLocal(reader)
		if err != nil {
			return nil, err
		}
		locals[i] = local
	}

	var insts []instruction.Instruction
	for {
		inst, err := ReadInstruction(reader)
		if err != nil {
			return nil, err
		}
		insts = append(insts, inst)
		_, ok := inst.(instruction.End)
		if ok {
			break
		}
	}

	return &Code{
		Size:       size,
		Locals:     locals,
		Expression: insts,
	}, nil
}

func ReadLocal(reader io.Reader) (Local, error) {
	types, err := ReadValueTypeVector(reader)
	if err != nil {
		return Local{}, err
	}
	return Local{
		ValueTypes: types,
	}, nil
}

func ReadInstruction(reader io.Reader) (instruction.Instruction, error) {
	opCode, err := ReadOpCode(reader)
	if err != nil {
		return nil, err
	}
	switch opCode {
	case opcode.End:
		return instruction.End{}, nil
	case opcode.LocalGet:
		index, err := ReadLebU128(reader)
		if err != nil {
			return nil, err
		}
		return instruction.LocalGet{
			Index: indicies.Local(index),
		}, nil
	case opcode.I32Add:
		return instruction.I32Add{}, nil
	}
	return nil, fmt.Errorf("invalid opcode %d", opCode)
}

func ReadOpCode(reader io.Reader) (opcode.Opcode, error) {
	b, err := ReadByte(reader)
	if err != nil {
		return 0, err
	}
	return opcode.Opcode(b), nil
}

func ReadUInt32(reader io.Reader) (uint32, error) {
	return read[uint32](reader)
}

func ReadUInt16(reader io.Reader) (uint16, error) {
	return read[uint16](reader)
}

func ReadByte(reader io.Reader) (byte, error) {
	return read[byte](reader)
}

func read[T any](reader io.Reader) (T, error) {
	var data T
	err := binary.Read(reader, binary.LittleEndian, &data)
	return data, err
}

func ReadLebU128(reader io.Reader) (uint32, error) {
	value, _, err := leb128.DecodeReader(reader)
	return value, err
}
