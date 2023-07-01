package types

type ExternType interface {
	Type
	externtype()
}

type ExternTypeImpl struct {
	TypeImpl
}

func (*ExternTypeImpl) externtype() {}
