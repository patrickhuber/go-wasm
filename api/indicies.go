package api

type Index interface {
	index()
}

type FuncIndex uint32

func (FuncIndex) index() {}

type TypeIndex uint32

func (TypeIndex) index() {}

type TableIndex uint32

func (TableIndex) index() {}

type MemoryIndex uint32

func (MemoryIndex) index() {}

type GlobalIndex uint32

func (GlobalIndex) index() {}

type ElementIndex uint32

func (ElementIndex) index() {}

type DataIndex uint32

func (DataIndex) index() {}

type LocalIndex uint32

func (LocalIndex) index() {}

type LabelIndex uint32

func (LabelIndex) index() {}
