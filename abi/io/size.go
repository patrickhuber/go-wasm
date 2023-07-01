package io

import (
	"math"

	"github.com/patrickhuber/go-wasm/abi/types"
)

const (
	SizeOfBool    = SizeOfU8
	SizeOfU8      = 1
	SizeOfU16     = 2
	SizeOfU32     = 4
	SizeOfU64     = 8
	SizeOfS8      = SizeOfU8
	SizeOfS16     = SizeOfU16
	SizeOfS32     = SizeOfU32
	SizeOfS64     = SizeOfU64
	SizeOfFloat32 = SizeOfU32
	SizeOfFloat64 = SizeOfU64
	SizeOfChar    = SizeOfU32
)

func Size(vt types.ValType) (uint32, error) {
	switch t := vt.(type) {
	case types.Bool:
		return 1, nil
	case types.U8:
		return SizeOfU8, nil
	case types.U16:
		return SizeOfU16, nil
	case types.U32:
		return SizeOfU64, nil
	case types.S8:
		return SizeOfS8, nil
	case types.S16:
		return SizeOfS16, nil
	case types.S32:
		return SizeOfS32, nil
	case types.Record:
		return SizeRecord(t)
	case types.Variant:
		return SizeVariant(t)
	case types.Flags:
		return SizeFlags(t)
	case types.Own:
		return 4, nil
	case types.Borrow:
		return 4, nil
	}
	return 0, nil
}

func SizeRecord(r types.Record) (uint32, error) {
	var s uint32 = 0
	for _, f := range r.Fields() {
		alignment, err := Alignment(f.Type)
		if err != nil {
			return 0, err
		}
		s, err = AlignTo(s, alignment)
		if err != nil {
			return 0, err
		}
		size, err := Size(f.Type)
		if err != nil {
			return 0, err
		}
		s += size
	}
	alignment, err := AlignmentRecord(r)
	if err != nil {
		return 0, err
	}
	return AlignTo(s, alignment)
}

func SizeVariant(v types.Variant) (uint32, error) {

	dt, err := DiscriminantType(v.Cases())
	if err != nil {
		return 0, err
	}

	s, err := Size(dt)
	if err != nil {
		return 0, err
	}

	alignment, err := MaxCaseAlignment(v.Cases())
	if err != nil {
		return 0, err
	}

	s, err = AlignTo(s, alignment)
	var cs uint32 = 0
	for _, c := range v.Cases() {
		if c.Type == nil {
			continue
		}
		size, err := Size(c.Type)
		if err != nil {
			return 0, err
		}
		cs = max(cs, size)
	}
	s += cs
	alignment, err = AlignmentVariant(v)
	return AlignTo(s, alignment)
}

func SizeFlags(f types.Flags) (uint32, error) {
	n := len(f.Labels())
	if n == 0 {
		return 0, nil
	}
	if n <= 8 {
		return 1, nil
	}
	if n <= 16 {
		return 2, nil
	}
	return 4 * NumI32Flags(f.Labels()), nil
}

func NumI32Flags(labels []string) uint32 {
	return uint32(math.Ceil(float64(len(labels)) / float64(32)))
}
