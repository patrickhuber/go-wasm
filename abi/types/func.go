package types

type Parameter struct {
	Name string
	Type ValType
}

type FuncType struct {
	Parameters []Parameter
	Results    []Parameter
}

func (ft *FuncType) ParamTypes() []ValType {
	return ft.extract(ft.Parameters)
}

func (ft *FuncType) ResultTypes() []ValType {
	return ft.extract(ft.Results)
}

func (ft *FuncType) extract(vec []Parameter) []ValType {
	var ret []ValType
	for _, t := range vec {
		ret = append(ret, t.Type)
	}
	return ret
}
