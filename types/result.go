package types

type Result interface {
	result()
	Types() []Value
}

type ResultImpl struct {
	types []Value
}

func (*ResultImpl) result() {}

func (r *ResultImpl) Types() []Value {
	return r.types
}
