package io

import (
	"fmt"
)

var ErrCast error = fmt.Errorf("error casting")

func NewCastError(v any, to string) error {
	from := fmt.Sprintf("%T", v)
	return fmt.Errorf("%v to %s : %w", from, to, ErrCast)
}
