package types

const (
	MaxStringByteLength uint32 = (1 << 31) - 1
)

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
