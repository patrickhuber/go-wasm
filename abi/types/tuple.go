package types

type Tuple interface {
	ValType
	Types() []ValType
	tuple()
}

type TupleImpl struct {
	ValTypeImpl
	types []ValType
}

func (*TupleImpl) tuple() {}

func (t *TupleImpl) Types() []ValType {
	return t.types
}

func NewTuple(types ...ValType) Tuple {
	return &TupleImpl{
		types: types,
	}
}
