package api

type Table struct {
	Limits    Limits
	Reference Reference
}

func (*Table) external() {}

type TableGet struct {
	Index TableIndex
}

func (*TableGet) instruction() {}

type TableSet struct {
	Index TableIndex
}

func (*TableSet) instruction() {}

type TableSize struct {
	Index TableIndex
}

func (*TableSize) instruction() {}

type TableGrow struct {
	Index TableIndex
}

func (*TableGrow) instruction() {}

type TableFill struct {
	Index TableIndex
}

func (*TableFill) instruction() {}

type TableCopy struct {
	Source      TableIndex
	Destination TableIndex
}

func (*TableCopy) instruction() {}

type TableInit struct {
	Destination TableIndex
	Source      ElementIndex
}

func (*TableInit) instruction() {}

type ElementDrop struct {
	Index ElementIndex
}

func (*ElementDrop) instruction() {}
