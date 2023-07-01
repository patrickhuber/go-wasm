package types

func NewS8() S8 {
	return &S8Impl{}
}

type S8 interface {
	ValType
	int8()
}

type S8Impl struct {
	ValTypeImpl
}

func (*S8Impl) int8() {}

type S16 interface {
	ValType
	int16()
}

type S16Impl struct {
	ValTypeImpl
}

func (*S16Impl) int16() {}

func NewS16() S16 {
	return &S16Impl{}
}

type S32 interface {
	ValType
	int32()
}

type S32Impl struct {
	ValTypeImpl
}

func (*S32Impl) int32() {}

func NewS32() S32 {
	return &S32Impl{}
}

type S64 interface {
	ValType
	int64()
}

type S64Impl struct {
	ValTypeImpl
}

func (*S64Impl) int64() {}

func NewS64() S64 {
	return &S64Impl{}
}
