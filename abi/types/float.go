package types

type Float32 interface {
	ValType
	float32()
}

type Float32Impl struct {
	ValTypeImpl
}

func (*Float32Impl) float32() {}

func NewFloat32() Float32 {
	return new(Float32Impl)
}

type Float64 interface {
	ValType
	float64()
}

type Float64Impl struct {
	ValTypeImpl
}

func (*Float64Impl) float64() {}

func NewFloat64() Float64 {
	return new(Float64Impl)
}
