package types

import "bytes"

// ReallocFunc defines a memory reallocation signature
type ReallocFunc func(originalPtr, originalSize, alignment, newSize uint32) (ptr uint32, err error)
type PostReturnFunc func()

type CanonicalOptions struct {
	Memory         bytes.Buffer
	StringEncoding StringEncoding
	Realloc        ReallocFunc
	PostReturn     PostReturnFunc
}
