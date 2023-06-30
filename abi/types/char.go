package types

type Char interface {
	ValType
	char()
}

type CharImpl struct {
	ValTypeImpl
}

func (*CharImpl) char() {}

func NewChar() Char {
	return new(CharImpl)
}
