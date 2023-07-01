package types

func NewU8() U8 {
	return &U8Impl{}
}

type U8 interface {
	ValType
	uint8()
}

type U8Impl struct {
	ValTypeImpl
}

func (*U8Impl) uint8() {}

type U16 interface {
	ValType
	uint16()
}

type U16Impl struct {
	ValTypeImpl
}

func (*U16Impl) uint16() {}

func NewU16() U16 {
	return &U16Impl{}
}

type U32 interface {
	ValType
	uint32()
}

type U32Impl struct {
	ValTypeImpl
}

func (*U32Impl) uint32() {}

func NewU32() U32 {
	return &U32Impl{}
}

type U64 interface {
	ValType
	uint64()
}

type U64Impl struct {
	ValTypeImpl
}

func (*U64Impl) uint64() {}

func NewU64() U64 {
	return &U64Impl{}
}
