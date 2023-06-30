package types

type String interface {
	ValType
	string()
}

type StringImpl struct {
	ValTypeImpl
}

func (*StringImpl) string() {}

func NewString() String {
	return new(StringImpl)
}
