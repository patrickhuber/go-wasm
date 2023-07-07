package types

type Table struct {
	Limits    Limits
	Reference Reference
}

func (*Table) external() {}
