package types

type ReallocFunc func(int, int, int, int) int
type PostReturnFunc func()

type CanonicalOptions struct {
	Memory         []byte
	StringEncoding StringEncoding
	Realloc        ReallocFunc
	PostReturn     PostReturnFunc
}
