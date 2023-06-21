package types

import "fmt"

type trap struct{}

func (t *trap) Error() string {
	return ""
}
func Trap() error {
	return &trap{}
}

func TrapWith(message string, args ...any) error {
	err := Trap()
	args = append(args, err)
	return fmt.Errorf(message+" %w", args...)
}
