package types

type Result interface {
	ValType
	Ok() ValType
	Error() ValType
	result()
}

type ResultImpl struct {
	ValTypeImpl
	ok  ValType
	err ValType
}

func (*ResultImpl) result() {}

func (r *ResultImpl) Ok() ValType {
	return r.ok
}

func (r *ResultImpl) Error() ValType {
	return r.err
}

func NewResult(ok ValType, err ValType) Result {
	return &ResultImpl{
		ok:  ok,
		err: err,
	}
}
