package values

type Vector interface {
	vector()
	Value
}

type V128Const struct {
	// Hi is the high order 64 bit part
	Hi uint64
	// Lo is the low order 64 bit part
	Lo uint64
}

func (*V128Const) vector() {}
func (*V128Const) value()  {}
