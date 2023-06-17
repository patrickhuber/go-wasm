package stack

// Push pushes the item t onto the stack (by appending) and returns a new slice
// with the item appended
func Push[T any](stack []T, t T) []T {
	return append(stack, t)
}

// Pop pops the item T off the stack by removing the last item and returns a new slice with
// the last item removed
// the popped item and a boolean true if the stack has items or false otherwise
func Pop[T any](stack []T) ([]T, T, bool) {
	var t T
	if len(stack) == 0 {
		return stack, t, false
	}
	t = stack[len(stack)-1]
	return stack[:len(stack)-1], t, true
}
