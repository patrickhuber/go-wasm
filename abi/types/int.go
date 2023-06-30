package types

func NewInt8() Int8 {
	return &Int8Impl{}
}

type Int8 interface {
	Type
	int8()
}

type Int8Impl struct {
	TypeImpl
}

func (*Int8Impl) int8() {}

type Int16 interface {
	Type
	int16()
}

type Int16Impl struct {
	TypeImpl
}

func (*Int16Impl) int16() {}

func NewInt16() Int16 {
	return &Int16Impl{}
}

type Int32 interface {
	Type
	int32()
}

type Int32Impl struct {
	TypeImpl
}

func (*Int32Impl) int32() {}

func NewInt32() Int32 {
	return &Int32Impl{}
}

type Int64 interface {
	Type
	int64()
}

type Int64Impl struct {
	TypeImpl
}

func (*Int64Impl) int64() {}

func NewInt64() Int64 {
	return &Int64Impl{}
}
