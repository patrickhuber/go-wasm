package to

func Pointer[T any](item T) *T {
	return &item
}
