package instructions

import (
	"github.com/patrickhuber/go-wasm/indicies"
	"github.com/patrickhuber/go-wasm/types"
)

type RefNull struct {
	ReferenceType types.Reference
}

func (*RefNull) instruction() {}

type RefIsNull struct {
}

func (*RefIsNull) instruction() {}

type RefFunc struct {
	FunctionIndex indicies.Function
}

func (*RefFunc) instruction() {}
