package types

type Enum interface {
	ValType
	Labels() []string
	enum()
}

type EnumImpl struct {
	ValTypeImpl
	labels []string
}

func (*EnumImpl) enum() {}

func (e *EnumImpl) Labels() []string {
	return e.labels
}

func NewEnum(labels ...string) Enum {
	return &EnumImpl{
		labels: labels,
	}
}
