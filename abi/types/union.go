package types

type Union interface {
	ValType
	Types() []ValType
	union()
}

type UnionImpl struct {
	ValTypeImpl
	types []ValType
}

func (*UnionImpl) union() {}

func (u *UnionImpl) Types() []ValType {
	return u.types
}

func NewUnion(types ...ValType) Union {
	return &UnionImpl{
		types: types,
	}
}
