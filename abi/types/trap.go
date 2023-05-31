package types

type Trap struct{}

func (t *Trap) Error() string {
	return ""
}

func TrapIf(condition bool) error {
	if condition {
		return &Trap{}
	}
	return nil
}
