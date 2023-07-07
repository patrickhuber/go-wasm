package types

type Number interface {
	number()
}

type I32 uint32

func (I32) number() {}
func (I32) value()  {}

type I64 uint64

func (I64) number() {}
func (I64) value()  {}

type F32 float32

func (F32) number() {}
func (F32) value()  {}

type F64 float64

func (F64) number() {}
func (F64) value()  {}
