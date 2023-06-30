package types

func NewUInt8() UInt8 {
	return &UInt8Impl{}
}

type UInt8 interface {
	ValType
	uint8()
}

type UInt8Impl struct {
	ValTypeImpl
}

func (*UInt8Impl) uint8() {}

type UInt16 interface {
	ValType
	uint16()
}

type UInt16Impl struct {
	ValTypeImpl
}

func (*UInt16Impl) uint16() {}

func NewUInt16() UInt16 {
	return &UInt16Impl{}
}

type UInt32 interface {
	ValType
	uint32()
}

type UInt32Impl struct {
	ValTypeImpl
}

func (*UInt32Impl) uint32() {}

func NewUInt32() UInt32 {
	return &UInt32Impl{}
}

type UInt64 interface {
	ValType
	uint64()
}

type UInt64Impl struct {
	ValTypeImpl
}

func (*UInt64Impl) uint64() {}

func NewUInt64() UInt64 {
	return &UInt64Impl{}
}
