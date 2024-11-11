package api

type Drop struct{}

func (*Drop) instruction() {}

type Select struct {
	Types []ValType
}

func (*Select) instruction() {}
