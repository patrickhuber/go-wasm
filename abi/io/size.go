package io

import (
	"math"

	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/abi/types"
)

const (
	SizeOfBool           = SizeOfU8
	SizeOfU8      uint32 = 1
	SizeOfU16     uint32 = 2
	SizeOfU32     uint32 = 4
	SizeOfU64     uint32 = 8
	SizeOfS8             = SizeOfU8
	SizeOfS16            = SizeOfU16
	SizeOfS32            = SizeOfU32
	SizeOfS64            = SizeOfU64
	SizeOfFloat32        = SizeOfU32
	SizeOfFloat64        = SizeOfU64
	SizeOfChar           = SizeOfU32
)

func Size(vt types.ValType) Result[uint32] {
	vt = Despecialize(vt)
	switch t := vt.(type) {
	case types.U8, types.S8, types.Bool:
		return result.Ok(SizeOfU8)
	case types.U16, types.S16:
		return result.Ok(SizeOfU16)
	case types.U32, types.S32, types.Float32:
		return result.Ok(SizeOfU32)
	case types.U64, types.S64, types.Float64:
		return result.Ok(SizeOfU64)
	case types.Char:
		return result.Ok(SizeOfChar)
	case types.List, types.String:
		return result.Ok(uint32(8))
	case types.Record:
		return SizeRecord(t)
	case types.Variant:
		return SizeVariant(t)
	case types.Flags:
		return SizeFlags(t)
	case types.Own, types.Borrow:
		return result.Ok(uint32(4))
	}
	return result.Errorf[uint32]("size: unable to match type %T", vt)
}

func SizeRecord(r types.Record) (res Result[uint32]) {
	defer handle.Error(&res)
	var s uint32 = 0
	for _, f := range r.Fields() {
		alignment := Alignment(f.Type).Unwrap()
		s = AlignTo(s, alignment)
		size := Size(f.Type).Unwrap()
		s += size
	}
	alignment := AlignmentRecord(r).Unwrap()
	return result.Ok(AlignTo(s, alignment))
}

func SizeVariant(v types.Variant) (res Result[uint32]) {
	defer handle.Error(&res)
	dt := DiscriminantType(v.Cases()).Unwrap()
	s := Size(dt).Unwrap()
	alignment := MaxCaseAlignment(v.Cases()).Unwrap()
	s = AlignTo(s, alignment)
	var cs uint32 = 0
	for _, c := range v.Cases() {
		if c.Type == nil {
			continue
		}
		size := Size(c.Type).Unwrap()
		cs = max(cs, size)
	}
	s += cs
	alignment = AlignmentVariant(v).Unwrap()
	return result.Ok(AlignTo(s, alignment))
}

func SizeFlags(f types.Flags) Result[uint32] {
	n := len(f.Labels())
	if n == 0 {
		return result.Ok(uint32(0))
	}
	if n <= 8 {
		return result.Ok(uint32(1))
	}
	if n <= 16 {
		return result.Ok(uint32(2))
	}
	return result.Ok(4 * NumI32Flags(f.Labels()))
}

func NumI32Flags(labels []string) uint32 {
	return uint32(math.Ceil(float64(len(labels)) / float64(32)))
}
