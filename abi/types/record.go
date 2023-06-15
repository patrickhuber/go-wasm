package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Field struct {
	Label string
	Type  ValType
}

type Record struct {
	Fields []Field
}

func (*Record) Kind() kind.Kind {
	return kind.Record
}

func (r *Record) Size() (uint32, error) {
	var s uint32 = 0
	return s, nil
}

func (r *Record) Alignment() (uint32, error) {
	var a uint32 = 1
	for _, f := range r.Fields {
		alignment, err := f.Type.Alignment()
		if err != nil {
			return 0, err
		}
		a = max(a, alignment)
	}
	return a, nil
}

func (r *Record) Despecialize() ValType {
	return r
}

func (r *Record) Flatten() ([]kind.Kind, error) {
	var flat []kind.Kind
	for _, f := range r.Fields {
		flattened, err := f.Type.Flatten()
		if err != nil {
			return nil, err
		}
		flat = append(flat, flattened...)
	}
	return flat, nil
}
