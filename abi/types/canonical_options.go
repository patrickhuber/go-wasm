package types

import (
	"bytes"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-wasm/encoding"
)

// ReallocFunc defines a memory reallocation signature
type ReallocFunc func(originalPtr, originalSize, alignment, newSize uint32) (ptr types.Result[uint32])
type PostReturnFunc func()

type CanonicalOptions struct {
	Memory         *bytes.Buffer
	StringEncoding encoding.Encoding
	Realloc        ReallocFunc
	PostReturn     PostReturnFunc
}
