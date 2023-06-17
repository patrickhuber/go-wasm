package to

// Pointer creates a pointer to the struct type T
func Pointer[T any](item T) *T {
	return &item
}
