package api

type Reference interface {
	reference()
}

type FunctionReference struct {
}

func (*FunctionReference) reference() {}
func (*FunctionReference) valType()   {}

type ExternalReference struct {
}

func (*ExternalReference) reference() {}
func (*ExternalReference) valType()   {}
