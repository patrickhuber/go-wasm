package types

type Bool interface {
	ValType
	bool()
}

type BoolImpl struct {
	ValTypeImpl
}

func (*BoolImpl) bool() {}

func NewBool() Bool {
	return new(BoolImpl)
}
