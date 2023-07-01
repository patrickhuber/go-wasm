package types

type Option interface {
	ValType
	Type() ValType
	option()
}
type OptionImpl struct {
	ValTypeImpl
	val ValType
}

func (*OptionImpl) option() {}

func (o *OptionImpl) Type() ValType {
	return o.val
}

func NewOption(val ValType) Option {
	return &OptionImpl{val: val}
}
