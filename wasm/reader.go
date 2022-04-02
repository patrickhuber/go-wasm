package wasm

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/patrickhuber/go-wasm/to"
)

type Reader interface {
	Read() (*Module, error)
	ReadLebU128() (uint32, error)
}

func NewReader(r io.Reader) Reader {
	return &reader{
		reader: r,
	}
}

type reader struct {
	reader io.Reader
}

func (r *reader) Read() (*Module, error) {
	module, err := r.readHeader()
	if err != nil {
		return nil, err
	}

	module.Sections, err = r.readSections()
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (r *reader) readHeader() (*Module, error) {
	module := &Module{}
	err := binary.Read(r.reader, binary.BigEndian, &module.Magic)
	if err != nil {
		return nil, err
	}
	err = binary.Read(r.reader, binary.LittleEndian, &module.Version)
	if err != nil {
		return nil, err
	}
	return module, err
}

func (r *reader) readSections() ([]Section, error) {
	var sections []Section
	for {
		section, err := r.readSection()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch section.ID {
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

		sections = append(sections, *section)
	}
	return sections, nil
}

func (r *reader) readSection() (*Section, error) {
	section := &Section{}
	err := binary.Read(r.reader, binary.LittleEndian, &section.ID)
	if err != nil {
		return nil, err
	}
	section.Size, err = r.ReadLebU128()
	if err != nil {
		return nil, err
	}
	return section, nil
}

func (r *reader) readTypeSection() (*TypeSection, error) {
	typeSection := &TypeSection{}
	count, err := r.ReadLebU128()
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

func (r *reader) readFunctionSection() (*FunctionSection, error) {
	funcSection := &FunctionSection{}
	count, err := r.ReadLebU128()
	if err != nil {
		return nil, err
	}
	for i := uint32(0); i < count; i++ {
		index, err := r.ReadLebU128()
		if err != nil {
			return nil, err
		}
		funcSection.Types = append(funcSection.Types, index)
	}
	return funcSection, nil
}

func (r *reader) readCodeSection() (*CodeSection, error) {
	codeSection := &CodeSection{}
	count, err := r.ReadLebU128()
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

func (r *reader) readCode() (*Code, error) {
	code := &Code{}
	size, err := r.ReadLebU128()
	if err != nil {
		return nil, err
	}
	code.Size = size

	localCount, err := r.ReadLebU128()
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

func (r *reader) readLocal() (*Locals, error) {
	local := &Locals{}
	t, err := r.readValueType()
	if err != nil {
		return nil, err
	}
	local.Type = t
	return local, nil
}

func (r *reader) readInstruction() (*Instruction, error) {
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

func (r *reader) readLocalInstruction() (*LocalInstruction, error) {
	index, err := r.ReadLebU128()
	if err != nil {
		return nil, err
	}
	return &LocalInstruction{
		Index: to.Pointer(index),
	}, nil
}

func (r *reader) readOpCode() (OpCode, error) {
	b, err := r.readByte()
	return OpCode(b), err
}

func (r *reader) readFuncType() (*Type, error) {
	b, err := r.readByte()
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

func (r *reader) readResultType() (*ResultType, error) {
	result := &ResultType{}
	size, err := r.ReadLebU128()
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

func (r *reader) readValueType() (*ValueType, error) {
	b, err := r.readByte()
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

func (r *reader) readByte() (byte, error) {
	b := make([]byte, 1)
	_, err := r.reader.Read(b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *reader) ReadLebU128() (uint32, error) {
	buf := make([]byte, 1)
	var val uint32
	shift := 0
	for {
		_, err := r.reader.Read(buf)
		if err != nil {
			return 0, err
		}
		val |= (uint32(buf[0]&0b_0111_1111) << shift)
		if buf[0]&0b_1000_0000 == 0 {
			break
		}
		shift += 7
	}
	return val, nil
}
