package values

import "github.com/patrickhuber/go-wasm/abi/kind"

type Int32Value struct {
	value int32
}

func (i *Int32Value) Kind() kind.Kind {
	return kind.S32
}

func (i *Int32Value) Value() any {
	return i.value
}

type Int64Value struct {
	value int64
}

func (*Int64Value) Kind() kind.Kind {
	return kind.S64
}

func (i *Int64Value) Value() any {
	return i.value
}
