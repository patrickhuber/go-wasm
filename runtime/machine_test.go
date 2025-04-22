package runtime_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/go-wasm/api"
	"github.com/patrickhuber/go-wasm/runtime"
	"github.com/patrickhuber/go-wasm/values"
)

func TestMachine(t *testing.T) {
	type test struct {
		name            string
		module          *api.Module
		externals       []values.Value
		expectedErr     error
		expectedReturns []values.Value
	}
	tests := []test{
		{
			name:            "empty",
			module:          &api.Module{},
			externals:       nil,
			expectedErr:     nil,
			expectedReturns: nil,
		},
		{
			name:            "i32_add",
			module:          &api.Module{},
			externals:       []values.Value{values.I32Const(1), values.I32Const(2)},
			expectedErr:     nil,
			expectedReturns: []values.Value{values.I32Const(3)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			machine := runtime.NewMachine()
			returns, err := machine.Execute(test.module, test.externals)
			if test.expectedErr != nil && err == nil {
				t.Errorf("expected error but found nil")
			}
			if test.expectedErr == nil && err != nil {
				t.Error(err)
			}
			if len(returns) != len(test.expectedReturns) {
				t.Errorf("expected %d returns but found %d", len(test.expectedReturns), len(returns))
			}
			for i, r := range returns {
				expectedReturn := test.expectedReturns[i]
				if reflect.DeepEqual(r, expectedReturn) {
					continue
				}
				t.Error("expected returns to be equal")
			}
		})
	}

}
