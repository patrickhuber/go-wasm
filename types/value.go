package types

type Value interface {
	value()
}

type ValueImpl struct {
}

func (*ValueImpl) value() {}
