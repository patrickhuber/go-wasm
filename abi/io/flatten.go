package io

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/abi/types"
)

func FlattenType(t types.ValType) ([]kind.Kind, error) {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool:
		return []kind.Kind{kind.U32}, nil
	case types.U8:
		return []kind.Kind{kind.U32}, nil
	case types.U16:
		return []kind.Kind{kind.U32}, nil
	case types.U32:
		return []kind.Kind{kind.U32}, nil
	case types.U64:
		return []kind.Kind{kind.U64}, nil
	case types.S8:
		return []kind.Kind{kind.U32}, nil
	case types.S16:
		return []kind.Kind{kind.U32}, nil
	case types.S32:
		return []kind.Kind{kind.U32}, nil
	case types.S64:
		return []kind.Kind{kind.U64}, nil
	case types.Float32:
		return []kind.Kind{kind.Float32}, nil
	case types.Float64:
		return []kind.Kind{kind.Float64}, nil
	case types.Char:
		return []kind.Kind{kind.U32}, nil
	case types.String:
		return []kind.Kind{kind.U32, kind.U32}, nil
	case types.List:
		return []kind.Kind{kind.U32, kind.U32}, nil
	case types.Record:
		return FlattenRecord(vt)
	case types.Variant:
		return FlattenVariant(vt)
	case types.Flags:
		flat := []kind.Kind{}
		n := NumI32Flags(vt.Labels())
		for i := uint32(0); i < n; i++ {
			flat = append(flat, kind.U32)
		}
		return flat, nil
	case types.Own:
		return []kind.Kind{kind.U32}, nil
	case types.Borrow:
		return []kind.Kind{kind.U32}, nil
	}
	return nil, fmt.Errorf("flatten_type: unable to match type %T", t)
}

func FlattenRecord(r types.Record) ([]kind.Kind, error) {
	flat := []kind.Kind{}
	for _, f := range r.Fields() {
		flattened, err := FlattenType(f.Type)
		if err != nil {
			return nil, err
		}
		flat = append(flat, flattened...)
	}
	return flat, nil
}

func FlattenVariant(v types.Variant) ([]kind.Kind, error) {
	var flat []kind.Kind
	for _, c := range v.Cases() {
		if c.Type == nil {
			continue
		}
		flattened, err := FlattenType(c.Type)
		if err != nil {
			return nil, err
		}
		for i, ft := range flattened {
			if i < len(flat) {
				flat[i] = join(flat[i], ft)
			} else {
				flat = append(flat, ft)
			}
		}
	}
	dt, err := DiscriminantType(v.Cases())
	if err != nil {
		return nil, err
	}

	flattened, err := FlattenType(dt)
	if err != nil {
		return nil, err
	}
	return append(flattened, flat...), nil
}

func join(a kind.Kind, b kind.Kind) kind.Kind {
	if a == b {
		return a
	}
	switch {
	case a == kind.U32 && b == kind.Float32:
		return kind.U32
	case a == kind.Float32 && b == kind.U32:
		return kind.U32
	default:
		return kind.U64
	}
}
