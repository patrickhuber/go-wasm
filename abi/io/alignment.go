package io

import (
	"math"

	. "github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/abi/types"
)

func AlignTo(ptr, alignment uint32) uint32 {
	return uint32(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment))
}

func Alignment(t types.ValType) Result[uint32] {
	t = Despecialize(t)
	switch vt := t.(type) {
	case types.Bool, types.U8, types.S8:
		return result.Ok[uint32](1)
	case types.U16, types.S16:
		return result.Ok[uint32](2)
	case types.U32, types.S32, types.Float32:
		return result.Ok[uint32](4)
	case types.U64, types.S64, types.Float64:
		return result.Ok[uint32](8)
	case types.Char, types.String, types.List:
		return result.Ok[uint32](4)
	case types.Record:
		return AlignmentRecord(vt)
	case types.Variant:
		return AlignmentVariant(vt)
	case types.Flags:
		return AlignmentFlags(vt)
	case types.Own, types.Borrow:
		return result.Ok[uint32](4)
	}
	return result.Error[uint32](types.TrapWith("Alignment: unable to align type %T", t))
}

func AlignmentRecord(r types.Record) (res Result[uint32]) {
	defer handle.Error(&res)
	var a uint32 = 1
	for _, field := range r.Fields() {
		alignment := Alignment(field.Type).Unwrap()
		a = max(a, alignment)
	}
	return result.Ok(a)
}

func AlignmentVariant(v types.Variant) (res Result[uint32]) {
	defer handle.Error(&res)
	dt := DiscriminantType(v.Cases()).Unwrap()
	alignment := Alignment(dt).Unwrap()
	maxAlignment := MaxCaseAlignment(v.Cases()).Unwrap()
	return result.Ok(max(alignment, maxAlignment))
}

func AlignmentFlags(f types.Flags) Result[uint32] {
	n := len(f.Labels())
	if n <= 8 {
		return result.Ok(uint32(1))
	}
	if n <= 16 {
		return result.Ok(uint32(2))
	}
	return result.Ok(uint32(4))
}

func MaxCaseAlignment(cases []types.Case) (res Result[uint32]) {
	defer handle.Error(&res)
	var a uint32 = 1
	for _, c := range cases {
		if c.Type == nil {
			continue
		}
		alignment := Alignment(c.Type).Unwrap()
		a = max(a, alignment)
	}
	return result.Ok(a)
}

func DiscriminantType(cases []types.Case) Result[types.ValType] {
	n := len(cases)
	if n > (1 << 32) {
		return result.Errorf[types.ValType]("case length %d exceeds max %d", n, (1 << 32))
	}
	switch uint64(math.Ceil(math.Log2(float64(n)) / 8)) {
	case 0:
		return result.Ok[types.ValType](types.NewU8())
	case 1:
		return result.Ok[types.ValType](types.NewU8())
	case 2:
		return result.Ok[types.ValType](types.NewU16())
	case 3:
		return result.Ok[types.ValType](types.NewU32())
	}
	return result.Errorf[types.ValType]("DiscriminateType unable to match math.ceil(math.log2(%d)/8) to [0,1,2,3]", n)
}
