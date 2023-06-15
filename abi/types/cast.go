package types

import (
	"fmt"
)

var ErrCast error = fmt.Errorf("error casting")

func NewCastError(v any, to string) error {
	from := fmt.Sprintf("%T", v)
	return fmt.Errorf("%w : %v to %s", ErrCast, from, to)
}
