package types

import (
	"fmt"
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

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

func (v *Variant) Size() (uint32, error) {
	var s uint32 = 0
	dt, err := v.DiscriminantType()
	if err != nil {
		return 0, err
	}
	if dt != nil {
		s, err = dt.Size()
		if err != nil {
			return 0, err
		}
	}
	maxCaseAlignment, err := v.MaxCaseAlignment()
	if err != nil {
		return 0, err
	}
	s = AlignTo(s, maxCaseAlignment)
	var cs uint32 = 0
	for _, c := range v.Cases {
		if c.Type != nil {
			size, err := c.Type.Size()
			if err != nil {
				return 0, err
			}
			cs = max(cs, size)
		}
	}
	s += cs
	alignment, err := v.Alignment()
	if err != nil {
		return 0, err
	}
	return AlignTo(s, alignment), nil
}

func (v *Variant) Alignment() (uint32, error) {
	maxCaseAlignment, err := v.MaxCaseAlignment()
	if err != nil {
		return 0, err
	}
	var alignment uint32 = 0
	dt, err := v.DiscriminantType()
	if err != nil {
		return 0, err
	}
	if dt != nil {
		alignment, err = dt.Alignment()
		if err != nil {
			return 0, err
		}
	}
	return max(maxCaseAlignment, alignment), nil
}

func (v *Variant) Despecialize() ValType {
	return v
}

func (v *Variant) DiscriminantType() (ValType, error) {
	n := len(v.Cases)
	if n < 0 || n > (1<<32) {
		return nil, fmt.Errorf("case length is out of bounds [0..32] found %d", n)
	}
	ceil := uint64(math.Ceil(math.Log2(float64(n)) / 8))
	switch ceil {
	case 0:
		var u8 U8
		return u8, nil
	case 1:
		var u8 U8
		return u8, nil
	case 2:
		var u16 U16
		return u16, nil
	case 3:
		var u32 U32
		return u32, nil
	}
	// this should be either checked or return an error
	return nil, fmt.Errorf("expected case 0-3 found %d", ceil)
}

func (v *Variant) Flatten() ([]kind.Kind, error) {
	var flat []kind.Kind
	for _, c := range v.Cases {
		if c.Type == nil {
			continue
		}
		flattened, err := c.Type.Flatten()
		if err != nil {
			return nil, err
		}
		for i, ft := range flattened {
			if i < len(flat) {
				flat[i] = join(flat[i], ft)
			} else {
				flat = append(flat, ft)
			}
		}
	}
	dt, err := v.DiscriminantType()
	if err != nil {
		return nil, err
	}

	flattened, err := dt.Flatten()
	if err != nil {
		return nil, err
	}
	return append(flat, flattened...), nil
}

func join(a kind.Kind, b kind.Kind) kind.Kind {
	if a == b {
		return a
	}
	switch {
	case a == kind.U32 && b == kind.Float32:
		return kind.U32
	case a == kind.Float32 && b == kind.U32:
		return kind.U32
	default:
		return kind.U64
	}
}

func (v *Variant) MaxCaseAlignment() (uint32, error) {
	var a uint32 = 1
	for _, c := range v.Cases {
		if c.Type != nil {
			alignment, err := c.Type.Alignment()
			if err != nil {
				return 0, err
			}
			a = max(a, alignment)
		}
	}
	return a, nil
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
