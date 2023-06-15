package values

import (
	"fmt"

	"github.com/patrickhuber/go-wasm/abi/kind"
)

func Zero(k kind.Kind) (Value, error) {
	switch k {
	case kind.U32:
		return U32(0), nil
	case kind.U64:
		return U64(0), nil
	case kind.Float32:
		return Float32(0), nil
	case kind.Float64:
		return Float64(0), nil
	default:
		return nil, fmt.Errorf("unrecognized type kind.%s", k)
	}
}
