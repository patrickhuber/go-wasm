package api

type ValType interface {
	valueType()
}

type NumType interface {
	numType()
}

type numType int

func (numType) numType()   {}
func (numType) valueType() {}

const (
	I32Type numType = 0
	I64Type numType = 1
	F32Type numType = 2
	F64Type numType = 3
)

type VecType interface {
	vecType()
}

type vecType int

func (vecType) vecType()   {}
func (vecType) valueType() {}

const V128Type vecType = 4

type RefType interface {
	refType()
}

type refType int

func (refType) valueType() {}
func (refType) refType()   {}

const (
	FuncRefType   refType = 5
	ExternRefType refType = 6
)
