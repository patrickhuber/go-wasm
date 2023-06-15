package types

import (
	"math"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

type ValType interface {
	Kind() kind.Kind
	Size() (uint32, error)
	Alignment() (uint32, error)
	Despecialize() ValType
	Flatten() ([]kind.Kind, error)
}

func AlignTo(ptr, alignment uint32) uint32 {
	fptr := float64(ptr)
	falignment := float64(alignment)
	return uint32(math.Ceil(fptr/falignment)) * alignment
}
