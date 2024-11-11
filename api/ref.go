package api

type RefNull struct {
	ReferenceType Reference
}

func (*RefNull) instruction() {}

type RefIsNull struct {
}

func (*RefIsNull) instruction() {}

type RefFunc struct {
	FunctionIndex FunctionIndex
}

func (*RefFunc) instruction() {}
