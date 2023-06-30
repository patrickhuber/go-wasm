package types

type List interface {
	ValType
	Type() ValType
	list()
}

type ListImpl struct {
	ValTypeImpl
	val ValType
}

func (*ListImpl) list() {}

func (l *ListImpl) Type() ValType {
	return l.val
}

func NewList(val ValType) List {
	return &ListImpl{
		val: val,
	}
}
