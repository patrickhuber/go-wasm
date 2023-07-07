package types

type Memory struct {
	Limits Limits
}

func (*Memory) external() {}
