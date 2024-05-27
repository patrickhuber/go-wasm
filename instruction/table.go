package instruction

import "github.com/patrickhuber/go-wasm/indicies"

type TableGet struct {
	Index indicies.Table
}

func (*TableGet) instruction() {}

type TableSet struct {
	Index indicies.Table
}

func (*TableSet) instruction() {}

type TableSize struct {
	Index indicies.Table
}

func (*TableSize) instruction() {}

type TableGrow struct {
	Index indicies.Table
}

func (*TableGrow) instruction() {}

type TableFill struct {
	Index indicies.Table
}

func (*TableFill) instruction() {}

type TableCopy struct {
	Source      indicies.Table
	Destination indicies.Table
}

func (*TableCopy) instruction() {}

type TableInit struct {
	Destination indicies.Table
	Source      indicies.Element
}

func (*TableInit) instruction() {}

type ElementDrop struct {
	Index indicies.Element
}

func (*ElementDrop) instruction() {}
