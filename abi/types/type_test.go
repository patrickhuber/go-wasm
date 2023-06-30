package types_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
)

func TestTakesObject(t *testing.T) {
	i8 := &types.Int8Impl{}
	TakesObject := func(o types.Object) {}
	TakesObject(i8)
}

func TestReturnsObject(t *testing.T) {

	ReturnsObject := func() types.Object {
		return &types.Int8Impl{}
	}
	obj := ReturnsObject()
	switch obj.(type) {
	case types.Int8:
	default:
		t.Fail()
	}
}

func TestTakesType(t *testing.T) {
	i8 := &types.Int8Impl{}
	TakesObject := func(t types.Type) {}
	TakesObject(i8)
}

func TestReturnsType(t *testing.T) {

	ReturnsObject := func() types.Type {
		return types.NewInt8()
	}
	obj := ReturnsObject()
	switch obj.(type) {
	case types.Int8:
	default:
		t.Fail()
	}
}
