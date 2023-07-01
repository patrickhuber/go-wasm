package types_test

import (
	"testing"

	"github.com/patrickhuber/go-wasm/abi/types"
)

func TestTakesObject(t *testing.T) {
	i8 := &types.S8Impl{}
	TakesObject := func(o types.Object) {}
	TakesObject(i8)
}

func TestReturnsObject(t *testing.T) {

	ReturnsObject := func() types.Object {
		return &types.S8Impl{}
	}
	obj := ReturnsObject()
	switch obj.(type) {
	case types.S8:
	default:
		t.Fail()
	}
}

func TestTakesType(t *testing.T) {
	i8 := &types.S8Impl{}
	TakesObject := func(t types.Type) {}
	TakesObject(i8)
}

func TestReturnsType(t *testing.T) {

	ReturnsObject := func() types.Type {
		return types.NewU8()
	}
	obj := ReturnsObject()
	switch typ := obj.(type) {
	case types.U8:
	default:
		t.Fatalf("unable to match type %T", typ)
	}
}
