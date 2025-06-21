package api

type ValType interface {
	valueType()
}

type NumType interface {
	numType()
}

type I32Type struct{}

func (*I32Type) valueType() {}
func (*I32Type) numType()   {}

type I64Type struct{}

func (*I64Type) valueType() {}
func (*I64Type) numType()   {}

type F32Type struct{}

func (*F32Type) valueType() {}
func (*F32Type) numType()   {}

type F64Type struct{}

func (*F64Type) valueType() {}
func (*F64Type) numType()   {}

type VecType interface {
	vecType()
}

type V128Type struct{}

func (*V128Type) valueType() {}
func (*V128Type) vecType()   {}

type RefType interface {
	refType()
}

type FuncRefType struct{}

func (*FuncRefType) valueType() {}
func (*FuncRefType) refType()   {}

type ExternRefType struct{}

func (*ExternRefType) valueType() {}
func (*ExternRefType) refType()   {}
