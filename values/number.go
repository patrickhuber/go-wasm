package values

type Number interface {
	number()
	Value
}

type I32Const uint32

func (I32Const) number() {}
func (I32Const) value()  {}

type I64Const uint64

func (I64Const) number() {}
func (I64Const) value()  {}

type F32Const float32

func (F32Const) number() {}
func (F32Const) value()  {}

type F64Const float64

func (F64Const) number() {}
func (F64Const) value()  {}
