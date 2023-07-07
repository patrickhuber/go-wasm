package types

type Function interface {
	Parameters() Result
	Results() Result
	function()
	External
}

type FunctionImpl struct {
	parameters Result
	results    Result
}

func (*FunctionImpl) function() {}
func (*FunctionImpl) external() {}

func (f *FunctionImpl) Parameters() Result {
	return f.parameters
}

func (f *FunctionImpl) Results() Result {
	return f.results
}
