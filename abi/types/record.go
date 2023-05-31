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

func (r *Record) Size() uint32 {
	var s uint32 = 0
	return s
}

func (r *Record) Alignment() uint32 {
	var a uint32 = 1
	for _, f := range r.Fields {
		a = max(a, f.Type.Alignment())
	}
	return a
}

func (r *Record) Despecialize() ValType {
	return r
}
