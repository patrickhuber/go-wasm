package types

type Parameter struct {
	Name string
	Type ValType
}

type FuncType interface {
	ExternType
	ParamTypes() []ValType
	ResultTypes() []ValType
	functype()
}

type FuncTypeImpl struct {
	ExternTypeImpl
	Parameters []Parameter
	Results    []Parameter
}

func NewFuncType(params []Parameter, results []Parameter) FuncType {
	return &FuncTypeImpl{
		Parameters: params,
		Results:    results,
	}
}
func (*FuncTypeImpl) functype() {}

func (ft *FuncTypeImpl) ParamTypes() []ValType {
	return ft.extract(ft.Parameters)
}

func (ft *FuncTypeImpl) ResultTypes() []ValType {
	return ft.extract(ft.Results)
}

func (ft *FuncTypeImpl) extract(vec []Parameter) []ValType {
	var ret []ValType
	for _, t := range vec {
		ret = append(ret, t.Type)
	}
	return ret
}
