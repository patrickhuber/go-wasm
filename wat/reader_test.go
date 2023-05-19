package wat_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-wasm/wat"
	"github.com/stretchr/testify/require"
)

type ReadTest struct {
	Text  string
	Types []wat.NodeType
}

func TestCanRead(t *testing.T) {
	tests := []ReadTest{
		{"(module)", []wat.NodeType{
			wat.ModuleNodeType,
			wat.EndNodeType}},
		{"( module )", []wat.NodeType{
			wat.ModuleNodeType,
			wat.EndNodeType,
		}},
		{"(module (memory 1) (func))", []wat.NodeType{
			wat.ModuleNodeType,
			wat.MemoryNodeType,
			wat.ArgumentNodeType,
			wat.EndNodeType,
			wat.FuncNodeType,
			wat.EndNodeType,
			wat.EndNodeType,
		}},
		{"(module (func $alias ))", []wat.NodeType{
			wat.ModuleNodeType,
			wat.FuncNodeType,
			wat.ArgumentNodeType,
			wat.EndNodeType,
			wat.EndNodeType,
		}},
		{`(module 
			(func (param i32) (param i32) (local i32) (result i64) get_local 0 
			   get_local 1 
			   get_local 2 
			) 
		 )`, []wat.NodeType{
			wat.ModuleNodeType,   // (module
			wat.FuncNodeType,     // (func
			wat.ParamNodeType,    // (param
			wat.ArgumentNodeType, // i32
			wat.EndNodeType,      // ) end param
			wat.ParamNodeType,    // (param
			wat.ArgumentNodeType, // i32
			wat.EndNodeType,      // ) end param
			wat.LocalNodeType,    // (local
			wat.ArgumentNodeType, // i32
			wat.EndNodeType,      // ) end local
			wat.ResultNodeType,   // (result
			wat.ArgumentNodeType, // i32
			wat.EndNodeType,      // ) end result
			wat.ArgumentNodeType, // get_local
			wat.ArgumentNodeType, // 0
			wat.ArgumentNodeType, // get_local
			wat.ArgumentNodeType, // 1
			wat.ArgumentNodeType, // get_local
			wat.ArgumentNodeType, // 2
			wat.EndNodeType,      // ) end func
			wat.EndNodeType,      // ) end module
		}},
		{
			`(module
				(func (export "add") (param i32 i32) (result i32)
				  local.get 0
				  local.get 1
				  i32.add))`,
			[]wat.NodeType{
				wat.ModuleNodeType,   // (module
				wat.FuncNodeType,     // (func
				wat.ExportNodeType,   // (export
				wat.ArgumentNodeType, // "add"
				wat.EndNodeType,      // ) end export
				wat.ParamNodeType,    // (param
				wat.ArgumentNodeType, // i32
				wat.ArgumentNodeType, // i32
				wat.EndNodeType,      // ) end param
				wat.ResultNodeType,   // (result
				wat.ArgumentNodeType, // i32
				wat.EndNodeType,      // ) end result
				wat.ArgumentNodeType, // local.get
				wat.ArgumentNodeType, // 0
				wat.ArgumentNodeType, // local.get
				wat.ArgumentNodeType, // 1
				wat.ArgumentNodeType, // i32.add
				wat.EndNodeType,      // ) end func
				wat.EndNodeType,      // ) end module
			},
		},
	}
	for _, test := range tests {
		reader := wat.NewReader(strings.NewReader(test.Text))
		i := 0
		for {
			next, err := reader.Next()
			if err != nil {
				require.FailNow(t, err.Error())
			}
			if !next {
				break
			}
			current := reader.Current()
			require.NotNil(t, current)
			require.Equal(t, current.Type(), test.Types[i])
			i++
		}
		require.True(t, i == len(test.Types), "node count mismatch")
	}
}
