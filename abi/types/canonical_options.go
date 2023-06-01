package types

type ReallocFunc func(uint32, uint32, uint32, uint32) uint32
type PostReturnFunc func()

type CanonicalOptions struct {
	Memory         []byte
	StringEncoding StringEncoding
	Realloc        ReallocFunc
	PostReturn     PostReturnFunc
}
