package types

type Flags interface {
	ValType
	flags()
	Labels() []string
}

type FlagsImpl struct {
	ValTypeImpl
	labels []string
}

func (*FlagsImpl) flags() {}

func (f *FlagsImpl) Labels() []string {
	return f.labels
}

func NewFlags(labels ...string) Flags {
	return &FlagsImpl{
		labels: labels,
	}
}
