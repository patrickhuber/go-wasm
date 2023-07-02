package types

type CoreExternType interface {
	Type
	coreexterntype()
}

type CoreExternTypeImpl struct {
	TypeImpl
}

func (*CoreExternTypeImpl) coreexterntype() {}
