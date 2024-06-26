package io

import (
	"fmt"
	"math"

	"github.com/patrickhuber/go-wasm/abi/types"
)

func AlignTo(ptr, alignment uint32) (uint32, error) {
	return uint32(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment)), nil
}

func Alignment(t types.ValType) (uint32, error) {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool:
		return 1, nil
	case types.U8:
		return 1, nil
	case types.U16:
		return 2, nil
	case types.U32:
		return 4, nil
	case types.U64:
		return 8, nil
	case types.S8:
		return 1, nil
	case types.S16:
		return 2, nil
	case types.S32:
		return 4, nil
	case types.S64:
		return 8, nil
	case types.F32:
		return 4, nil
	case types.F64:
		return 8, nil
	case types.Char:
		return 4, nil
	case types.String:
		return 4, nil
	case types.List:
		return 4, nil
	case types.Record:
		return AlignmentRecord(vt)
	case types.Variant:
		return AlignmentVariant(vt)
	case types.Flags:
		return AlignmentFlags(vt)
	case types.Own:
		return 4, nil
	case types.Borrow:
		return 4, nil
	}
	return 0, types.TrapWith("Alignment: unable to align type %T", t)
}

func AlignmentRecord(r types.Record) (uint32, error) {
	var a uint32 = 1
	for _, field := range r.Fields() {
		alignment, err := Alignment(field.Type)
		if err != nil {
			return 0, err
		}
		a = max(a, alignment)
	}
	return a, nil
}

func AlignmentVariant(v types.Variant) (uint32, error) {
	dt, err := DiscriminantType(v.Cases())
	if err != nil {
		return 0, err
	}
	alignment, err := Alignment(dt)
	if err != nil {
		return 0, nil
	}
	maxAlignment, err := MaxCaseAlignment(v.Cases())
	if err != nil {
		return 0, nil
	}
	return max(alignment, maxAlignment), nil
}

func AlignmentFlags(f types.Flags) (uint32, error) {
	n := len(f.Labels())
	if n <= 8 {
		return 1, nil
	}
	if n <= 16 {
		return 2, nil
	}
	return 4, nil
}

func MaxCaseAlignment(cases []types.Case) (uint32, error) {
	var a uint32 = 1
	for _, c := range cases {
		if c.Type == nil {
			continue
		}
		alignment, err := Alignment(c.Type)
		if err != nil {
			return 0, err
		}
		a = max(a, alignment)
	}
	return a, nil
}

func DiscriminantType(cases []types.Case) (types.ValType, error) {
	n := len(cases)
	// n must be positive integer
	if !(0 < n && n < math.MaxInt) {
		return nil, fmt.Errorf("case length %d exceeds max %d", n, math.MaxInt)
	}
	switch uint64(math.Ceil(math.Log2(float64(n)) / 8)) {
	case 0:
		return types.NewU8(), nil
	case 1:
		return types.NewU8(), nil
	case 2:
		return types.NewU16(), nil
	case 3:
		return types.NewU32(), nil
	}
	return nil, fmt.Errorf("DiscriminateType unable to match math.ceil(math.log2(%d)/8) to [0,1,2,3]", n)
}
