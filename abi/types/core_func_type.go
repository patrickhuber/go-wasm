package types

import "github.com/patrickhuber/go-wasm/abi/kind"

type CoreFuncType interface {
	CoreExternType
	corefunctype()
	Params() []kind.Kind
	Results() []kind.Kind
}

type CoreFuncTypeImpl struct {
	CoreExternTypeImpl
	params  []kind.Kind
	results []kind.Kind
}

func (*CoreFuncTypeImpl) corefunctype() {}

func (cft *CoreFuncTypeImpl) Params() []kind.Kind {
	return cft.params
}

func (cft *CoreFuncTypeImpl) Results() []kind.Kind {
	return cft.results
}

func NewCoreFuncType(params []kind.Kind, results []kind.Kind) CoreFuncType {
	return &CoreFuncTypeImpl{
		params:  params,
		results: results,
	}
}
