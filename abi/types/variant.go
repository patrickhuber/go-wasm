package types

type Variant interface {
	ValType
	Cases() []Case
	variant()
}

type VariantImpl struct {
	ValTypeImpl
	cases []Case
}

// Cases implements Variant.
func (v *VariantImpl) Cases() []Case {
	return v.cases
}

// variant implements Variant.
func (*VariantImpl) variant() {}

type Case struct {
	Label   string
	Type    ValType
	Refines *string
}

func NewCase(label string, valType ValType) Case {
	return Case{
		Label:   label,
		Type:    valType,
		Refines: nil,
	}
}

func NewCaseRefines(label string, valType ValType, refines string) Case {
	return Case{
		Label:   label,
		Type:    valType,
		Refines: &refines,
	}
}

func NewVariant(cases ...Case) Variant {
	return &VariantImpl{
		cases: cases,
	}
}
