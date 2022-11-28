package runtime

import (
	"encoding/binary"
	"math"
)

// https://webassembly.github.io/spec/core/exec/runtime.html#syntax-val
type Value struct {
	Number *Number
	Vec    *Vec
	Ref    *Ref
}

func (v *Value) Uint32() uint32 {
	return binary.BigEndian.Uint32(v.Number.Value)
}

func (v *Value) Uint64() uint64 {
	return binary.BigEndian.Uint64(v.Number.Value)
}

func (v *Value) Float32() float32 {
	asInt := v.Uint32()
	return math.Float32frombits(asInt)
}

func (v *Value) Float64() float64 {
	asInt := v.Uint64()
	return math.Float64frombits(asInt)
}

func (v *Value) SetUint32(value uint32) {
	binary.BigEndian.PutUint32(v.Number.Value, value)
}

func (v *Value) SetUint64(value uint64) {
	binary.BigEndian.PutUint64(v.Number.Value, value)
}

func (v *Value) SetFloat32(value float32) {
	bits := math.Float32bits(value)
	v.SetUint32(bits)
}

func (v *Value) SetFloat64(value float64) {
	bits := math.Float64bits(value)
	v.SetUint64(bits)
}
