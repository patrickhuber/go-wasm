package types

type ValType interface {
	Type
	valtype()
}

type ValTypeImpl struct {
	TypeImpl
}

func (*ValTypeImpl) valtype() {}
