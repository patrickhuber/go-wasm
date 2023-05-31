package types

import "golang.org/x/exp/constraints"

func max[T constraints.Ordered](left, right T) T {
	if left > right {
		return left
	}
	return right
}
