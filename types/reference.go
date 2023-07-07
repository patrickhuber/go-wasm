package types

type Reference interface {
	reference()
}

type FunctionReference struct {
	ValueImpl
}

func (*FunctionReference) reference() {}

type ExternalReference struct {
	ValueImpl
}

func (*ExternalReference) reference() {}
