package binary

import (
	"errors"
	"fmt"
	"io"
	"slices"

	"encoding/binary"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/leb128"
	"github.com/patrickhuber/go-wasm/opcode"
)

func Read(reader io.Reader) (*api.Document, error) {

	preamble, err := ReadPreamble(reader)
	if err != nil {
		return nil, err
	}

	var directive api.Directive
	switch preamble.Version {
	case ComponentVersion:
		directive, err = ReadComponent(reader)
		if err != nil {
			return nil, err
		}
	case ModuleVersion:
		directive, err = ReadModule(reader)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid version %d", preamble.Version)
	}

	return &api.Document{
		Preamble:  preamble,
		Directive: directive,
	}, nil
}

func ReadPreamble(reader io.Reader) (api.Preamble, error) {
	magic, err := ReadBytes(reader, len(Magic))
	if err != nil {
		return api.Preamble{}, err
	}

	if !slices.Equal(Magic, magic[0:]) {
		return api.Preamble{}, fmt.Errorf("expected magic %v found %v", Magic, magic)
	}

	version, err := ReadUInt16(reader)
	if err != nil {
		return api.Preamble{}, err
	}

	layer, err := ReadUInt16(reader)
	if err != nil {
		return api.Preamble{}, err
	}

	return api.Preamble{
		Version: version,
		Layer:   layer,
	}, nil
}

func ReadModule(reader io.Reader) (*api.Module, error) {
	module := &api.Module{}
	for {
		sectionID, size, err := ReadSectionHeader(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		switch sectionID {
		case TypeSectionID:
			types, err := ReadFuncTypes(size, reader)
			if err != nil {
				return nil, err
			}
			module.Types = types
		case FunctionSectionID:
			funcs, err := ReadFuncs(size, reader)
			if err != nil {
				return nil, err
			}
			module.Funcs = funcs
		case CodeSectionID:
			err := UpdateFuncsWithCode(module, size, reader)
			if err != nil {
				return nil, err
			}

		case ExportSectionID:
			exports, err := ReadExports(size, reader)
			if err != nil {
				return nil, err
			}
			module.Exports = exports
		default:
			// skip unknown sections
			data := make([]byte, size)
			if size > 0 {
				if _, err := io.ReadFull(reader, data); err != nil {
					return nil, fmt.Errorf("failed to read section %d: %w", sectionID, err)
				}
			}
		}
	}
	return module, nil
}

func ReadComponent(reader io.Reader) (*api.Component, error) {
	return &api.Component{}, nil
}

func ReadSectionHeader(reader io.Reader) (SectionID, uint32, error) {
	id, err := ReadByte(reader)
	if err != nil {
		return SectionID(0), 0, err
	}

	// read the size
	size, err := ReadLebU128(reader)
	if err != nil {
		return SectionID(0), 0, err
	}
	return SectionID(id), size, nil
}

func ReadFuncTypes(size uint32, reader io.Reader) ([]*api.FuncType, error) {
	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}
	types := make([]*api.FuncType, count)
	for i := uint32(0); i < count; i++ {
		t, err := ReadType(reader)
		if err != nil {
			return nil, err
		}
		types[i] = t
	}
	return types, nil
}

func ReadType(reader io.Reader) (*api.FuncType, error) {
	b, err := ReadByte(reader)
	if err != nil {
		return nil, err
	}
	if b != 0x60 {
		return nil, fmt.Errorf("expected byte 0x60 but found %b", b)
	}
	parameters, err := ReadResultType(reader)
	if err != nil {
		return nil, err
	}
	results, err := ReadResultType(reader)
	if err != nil {
		return nil, err
	}
	return &api.FuncType{
		Parameters: parameters,
		Returns:    results,
	}, nil
}

func ReadResultType(reader io.Reader) (api.ResultType, error) {
	valueTypeVector, err := ReadValueTypeVector(reader)
	if err != nil {
		return api.ResultType{}, err
	}
	return api.ResultType{
		Types: valueTypeVector,
	}, nil
}

func ReadValueTypeVector(reader io.Reader) ([]api.ValType, error) {
	// read the size of the vector
	size, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}
	valueTypes := make([]api.ValType, size)
	for i := uint32(0); i < size; i++ {
		vt, err := ReadValueType(reader)
		if err != nil {
			return nil, err
		}
		valueTypes[i] = vt
	}
	return valueTypes, nil
}

func ReadValueType(reader io.Reader) (api.ValType, error) {
	b, err := ReadByte(reader)
	if err != nil {
		return nil, err
	}
	switch ValType(b) {
	case I32:
		return &api.I32Type{}, nil
	case I64:
		return &api.I64Type{}, nil
	case F32:
		return &api.F32Type{}, nil
	case F64:
		return &api.F64Type{}, nil
	}
	return nil, fmt.Errorf("invalid ValueType found %b", b)
}

func ReadFuncs(size uint32, reader io.Reader) ([]*api.Func, error) {

	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	funcs := make([]*api.Func, count)
	for i := uint32(0); i < count; i++ {
		result, err := ReadFunc(reader)
		if err != nil {
			return nil, err
		}
		funcs[i] = result
	}

	return funcs, nil
}

func ReadFunc(reader io.Reader) (*api.Func, error) {
	index, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}
	return &api.Func{
		Type: api.TypeIndex(index),
	}, nil
}

func ReadExports(size uint32, reader io.Reader) ([]api.Export, error) {
	count, err := ReadLebU128(reader)
	if err != nil {
		return nil, err
	}

	exports := make([]api.Export, count)
	for i := uint32(0); i < count; i++ {
		export, err := ReadExport(reader)
		if err != nil {
			return nil, err
		}
		exports[i] = export
	}

	return exports, nil
}

func ReadExport(reader io.Reader) (api.Export, error) {
	var zero api.Export
	name, err := ReadString(reader)
	if err != nil {
		return zero, err
	}

	exportKind, err := ReadByte(reader)
	if err != nil {
		return zero, err
	}

	index, err := ReadLebU128(reader)
	if err != nil {
		return zero, err
	}

	switch ExportKind(exportKind) {
	case FuncExportKind:
		return api.Export{
			Name: name,
			Description: &api.FuncExportDescription{
				FuncIdx: api.FuncIndex(index),
			},
		}, nil
	default:
		return zero, fmt.Errorf("invalid export kind %d", exportKind)
	}
}

func ReadString(reader io.Reader) (string, error) {
	size, err := ReadLebU128(reader)
	if err != nil {
		return "", err
	}

	if size == 0 {
		return "", nil
	}

	buf := make([]byte, size)
	err = binary.Read(reader, binary.LittleEndian, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func ReadExpression(reader io.Reader) (*api.Expression, error) {
	var insts []api.Instruction
	for {
		inst, err := ReadInstruction(reader)
		if err != nil {
			return nil, err
		}
		insts = append(insts, inst)
		_, ok := inst.(api.End)
		if ok {
			break
		}
	}
	return &api.Expression{
		Instructions: insts,
	}, nil
}

func UpdateFuncsWithCode(module *api.Module, size uint32, reader io.Reader) error {
	index := 0
	for {
		if index >= len(module.Funcs) {
			return fmt.Errorf("code section has more functions than function section, index %d, funcs %d", index, len(module.Funcs))
		}
		if index < 0 {
			return fmt.Errorf("code section has negative index %d", index)
		}

		_, err := ReadLebU128(reader)
		if err != nil {
			return err
		}

		fn := module.Funcs[index]
		locals, err := ReadValueTypeVector(reader)
		if err != nil {
			return err
		}
		fn.Locals = locals
		body, err := ReadExpression(reader)
		if err != nil {
			return err
		}
		fn.Body = body
	}
}

func ReadInstruction(reader io.Reader) (api.Instruction, error) {
	opCode, err := ReadOpCode(reader)
	if err != nil {
		return nil, err
	}
	switch opCode {
	case opcode.End:
		return api.End{}, nil
	case opcode.LocalGet:
		index, err := ReadLebU128(reader)
		if err != nil {
			return nil, err
		}
		return api.LocalGet{
			Index: api.LocalIndex(index),
		}, nil
	case opcode.I32Add:
		return api.I32Add{}, nil
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

func ReadBytes(reader io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	err := binary.Read(reader, binary.LittleEndian, buf)
	return buf, err
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
