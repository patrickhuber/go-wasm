package types

type Type interface {
	typ()
}

type TypeImpl struct {
	Object
}

func (*TypeImpl) typ() {}
