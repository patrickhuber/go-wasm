package wasm

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/leb128"
	"github.com/patrickhuber/go-wasm/internal/to"
)

type SectionReader interface {
	Read() (*Section, error)
}

type sectionReader struct {
	layer  Layer
	reader *bufio.Reader
}

func NewSectionReader(reader *bufio.Reader, layer Layer) SectionReader {
	return &sectionReader{
		reader: reader,
		layer:  layer,
	}
}

func (r *sectionReader) Read() (*Section, error) {
	section, err := r.readSection()
	if err != nil {
		return nil, err
	}
	if section.Size <= 0 {
		return section, nil
	}
	switch section.ID {
	case CustomSectionType:
		section.Custom, err = r.readCustomSection(section.Size)
	case TypeSectionType:
		section.Type, err = r.readTypeSection()
	case FuncSectionType:
		section.Function, err = r.readFunctionSection()
	case CodeSectionType:
		section.Code, err = r.readCodeSection()
	default:
		return nil, fmt.Errorf("unrecognized section ID %d", section.ID)
	}

	if err != nil {
		return nil, err
	}
	return section, nil
}

func (r *sectionReader) readSection() (*Section, error) {
	section := &Section{}
	err := binary.Read(r.reader, binary.LittleEndian, &section.ID)
	if err != nil {
		return nil, err
	}
	section.Size, err = r.readLebU128()
	if err != nil {
		return nil, err
	}
	return section, nil
}

func (r *sectionReader) readCustomSection(size uint32) (*CustomSection, error) {

	name, read, err := ReadUtf8String(r.reader)
	if err != nil {
		return nil, err
	}

	limit := size - uint32(read)
	nameSection, err := r.readNameSection(limit)
	if err != nil {
		return nil, err
	}
	nameSection.Key = name
	return &CustomSection{
		Name: nameSection,
	}, nil
}

func (r *sectionReader) readNameSection(limit uint32) (*NameSection, error) {
	buf := make([]byte, limit)
	_, err := io.ReadFull(r.reader, buf)
	if err != nil {
		return nil, err
	}
	for index := 0; index < len(buf); {
		subSectionID := buf[index]
		index++

		subSectionSize, read, err := leb128.DecodeSlice(buf[index:])
		if err != nil {
			return nil, err
		}
		index += read

		switch SubSectionID(subSectionID) {
		case SubSectionName:
			limit := int(subSectionSize) + index
			value, read, err := DecodeUtf8String(buf[index:limit])
			if err != nil {
				return nil, err
			}
			index += read
			return &NameSection{
				ID:   SubSectionID(subSectionID),
				Name: &value,
			}, nil
		default:
			return nil, fmt.Errorf("invalid subsection id %d", subSectionID)
		}

	}
	return nil, fmt.Errorf("unable to read subsection")
}

func (r *sectionReader) readTypeSection() (*TypeSection, error) {
	typeSection := &TypeSection{}
	count, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	for i := uint32(0); i < count; i++ {
		t, err := r.readFuncType()
		if err != nil {
			return nil, err
		}
		typeSection.Types = append(typeSection.Types, *t)
	}
	return typeSection, nil
}

func (r *sectionReader) readFunctionSection() (*FunctionSection, error) {
	funcSection := &FunctionSection{}
	count, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	for i := uint32(0); i < count; i++ {
		index, err := r.readLebU128()
		if err != nil {
			return nil, err
		}
		funcSection.Types = append(funcSection.Types, index)
	}
	return funcSection, nil
}

func (r *sectionReader) readCodeSection() (*CodeSection, error) {
	codeSection := &CodeSection{}
	count, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	for i := uint32(0); i < count; i++ {
		code, err := r.readCode()
		if err != nil {
			return nil, err
		}
		codeSection.Codes = append(codeSection.Codes, *code)
	}
	return codeSection, nil
}

func (r *sectionReader) readCode() (*Code, error) {
	code := &Code{}
	size, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	code.Size = size

	localCount, err := r.readLebU128()
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(localCount); i++ {
		local, err := r.readLocal()
		if err != nil {
			return nil, err
		}
		code.Locals = append(code.Locals, *local)
	}

	for {
		instr, err := r.readInstruction()
		if err != nil {
			return nil, err
		}
		code.Expression = append(code.Expression, *instr)
		if instr.OpCode == End {
			break
		}
	}
	return code, nil
}

func (r *sectionReader) readLocal() (*Locals, error) {
	local := &Locals{}
	t, err := r.readValueType()
	if err != nil {
		return nil, err
	}
	local.Type = t
	return local, nil
}

func (r *sectionReader) readInstruction() (*Instruction, error) {
	instruction := &Instruction{}
	opcode, err := r.readOpCode()
	if err != nil {
		return nil, err
	}
	instruction.OpCode = opcode
	switch {
	case LocalGet <= opcode && opcode <= LocalTee:
		instruction.Local, err = r.readLocalInstruction()
	case I32Clz <= opcode && opcode <= I32Add:
		// these are just opcode, no immedate
		return instruction, nil
	}
	if err != nil {
		return nil, err
	}

	return instruction, err
}

func (r *sectionReader) readLocalInstruction() (*LocalInstruction, error) {
	index, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	return &LocalInstruction{
		Index: to.Pointer(index),
	}, nil
}

func (r *sectionReader) readOpCode() (OpCode, error) {
	b, err := r.reader.ReadByte()
	return OpCode(b), err
}

func (r *sectionReader) readFuncType() (*Type, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != 0x60 {
		return nil, fmt.Errorf("expected function type prefix of 0x60")
	}
	parameters, err := r.readResultType()
	if err != nil {
		return nil, err
	}
	results, err := r.readResultType()
	if err != nil {
		return nil, err
	}
	return &Type{
		Parameters: parameters,
		Results:    results,
	}, nil
}

func (r *sectionReader) readResultType() (*ResultType, error) {
	result := &ResultType{}
	size, err := r.readLebU128()
	if err != nil {
		return nil, err
	}
	result.Values = make([]*ValueType, size)
	for i := uint32(0); i < size; i++ {
		value, err := r.readValueType()
		if err != nil {
			return nil, err
		}
		result.Values[i] = value
	}
	return result, nil
}

func (r *sectionReader) readValueType() (*ValueType, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return nil, err
	}
	v := &ValueType{}
	switch b {
	case byte(FuncRef), byte(ExternRef):
		v.ReferenceType = to.Pointer(ReferenceType(b))
	case byte(I32), byte(I64), byte(F32), byte(F64):
		v.NumberType = to.Pointer(NumberType(b))
	}
	return v, nil
}

func (r *sectionReader) readLebU128() (uint32, error) {
	value, _, err := leb128.Decode(r.reader)
	return value, err
}
