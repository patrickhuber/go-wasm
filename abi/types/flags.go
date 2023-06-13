package types

import (
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

type Flags struct {
	Labels []string
}

func (*Flags) Kind() kind.Kind {
	return kind.Flags
}

func (f *Flags) Alignment() uint32 {
	n := uint32(len(f.Labels))
	switch {
	case n <= 8:
		return 1
	case n <= 16:
		return 2
	}
	return 4
}

func (f *Flags) Despecialize() ValType {
	return f
}

func (f *Flags) Size() uint32 {
	n := len(f.Labels)
	switch {
	case n == 0:
		return 0
	case n <= 8:
		return 1
	case n <= 16:
		return 2
	}
	return uint32(4) * f.NumI32Flags()
}

func (f *Flags) Flatten() []kind.Kind {
	flat := []kind.Kind{}
	n := f.NumI32Flags()
	for i := uint32(0); i < n; i++ {
		flat = append(flat, kind.S32)
	}
	return flat
}

func (f *Flags) NumI32Flags() uint32 {
	flen := float64(len(f.Labels))
	f32 := float64(32)
	fdiv := flen / f32
	return uint32(math.Ceil(fdiv))
}
