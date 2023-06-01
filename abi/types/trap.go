package types

type trap struct{}

func (t *trap) Error() string {
	return ""
}
func Trap() error {
	return &trap{}
}

func TrapIf(condition bool) error {
	if condition {
		return &trap{}
	}
	return nil
}
