package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type Case struct {
	Label   string
	Type    ValType
	Refines *string
}

type Variant struct {
	Cases []Case
}

func (*Variant) Kind() kind.Kind {
	return kind.Variant
}

func (v *Variant) Size() uint32 {
	var s uint32 = 0
	// check nil
	if dt := v.DiscriminantType(); dt != nil {
		s = dt.Size()
	}
	s = AlignTo(s, v.MaxCaseAlignment())
	var cs uint32 = 0
	for _, c := range v.Cases {
		if c.Type != nil {
			cs = max(cs, c.Type.Size())
		}
	}
	s += cs
	return AlignTo(s, v.Alignment())
}

func (v *Variant) Alignment() uint32 {
	maxCaseAlignment := v.MaxCaseAlignment()
	var alignment uint32 = 0
	if dt := v.DiscriminantType(); dt != nil {
		alignment = dt.Alignment()
	}
	return max(maxCaseAlignment, alignment)
}

func (v *Variant) DiscriminantType() ValType {
	n := len(v.Cases)
	switch n {
	case 0:
		var u8 U8
		return u8
	case 1:
		var u8 U8
		return u8
	case 2:
		var u16 U16
		return u16
	case 3:
		var u32 U32
		return u32
	}
	// this should be either checked or return an error
	return nil
}

func (v *Variant) MaxCaseAlignment() uint32 {
	var a uint32 = 1
	for _, c := range v.Cases {
		if c.Type != nil {
			a = max(a, c.Type.Alignment())
		}
	}
	return a
}

func (v *Variant) CaseLabelWithRefinements(c Case) string {
	label := c.Label
	for c.Refines != nil {
		ind := v.findCaseIndex(*c.Refines)
		if ind < 0 {
			break
		}
		c = v.Cases[ind]
		label += "|" + c.Label
	}
	return label
}

func (v *Variant) findCaseIndex(label string) int {
	for i, c := range v.Cases {
		if c.Label == label {
			return i
		}
	}
	return -1
}
