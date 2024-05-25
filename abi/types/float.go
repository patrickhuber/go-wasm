package types

type F32 interface {
	ValType
	float32()
}

type F32Impl struct {
	ValTypeImpl
}

func (*F32Impl) float32() {}

func NewF32() F32 {
	return new(F32Impl)
}

type F64 interface {
	ValType
	float64()
}

type F64Impl struct {
	ValTypeImpl
}

func (*F64Impl) float64() {}

func NewF64() F64 {
	return new(F64Impl)
}
