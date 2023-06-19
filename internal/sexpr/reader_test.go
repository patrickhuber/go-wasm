package sexpr_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-wasm/internal/sexpr"
	"github.com/stretchr/testify/require"
)

type ReadTest struct {
	Text  string
	Types []sexpr.NodeType
}

func TestCanRead(t *testing.T) {
	tests := []ReadTest{
		{"(module)", []sexpr.NodeType{
			sexpr.ModuleNodeType,
			sexpr.EndNodeType}},
		{"( module )", []sexpr.NodeType{
			sexpr.ModuleNodeType,
			sexpr.EndNodeType,
		}},
		{"(module (memory 1) (func))", []sexpr.NodeType{
			sexpr.ModuleNodeType,
			sexpr.MemoryNodeType,
			sexpr.ArgumentNodeType,
			sexpr.EndNodeType,
			sexpr.FuncNodeType,
			sexpr.EndNodeType,
			sexpr.EndNodeType,
		}},
		{"(module (func $alias ))", []sexpr.NodeType{
			sexpr.ModuleNodeType,
			sexpr.FuncNodeType,
			sexpr.ArgumentNodeType,
			sexpr.EndNodeType,
			sexpr.EndNodeType,
		}},
		{`(module 
			(func (param i32) (param i32) (local i32) (result i64) get_local 0 
			   get_local 1 
			   get_local 2 
			) 
		 )`, []sexpr.NodeType{
			sexpr.ModuleNodeType,   // (module
			sexpr.FuncNodeType,     // (func
			sexpr.ParamNodeType,    // (param
			sexpr.ArgumentNodeType, // i32
			sexpr.EndNodeType,      // ) end param
			sexpr.ParamNodeType,    // (param
			sexpr.ArgumentNodeType, // i32
			sexpr.EndNodeType,      // ) end param
			sexpr.LocalNodeType,    // (local
			sexpr.ArgumentNodeType, // i32
			sexpr.EndNodeType,      // ) end local
			sexpr.ResultNodeType,   // (result
			sexpr.ArgumentNodeType, // i32
			sexpr.EndNodeType,      // ) end result
			sexpr.ArgumentNodeType, // get_local
			sexpr.ArgumentNodeType, // 0
			sexpr.ArgumentNodeType, // get_local
			sexpr.ArgumentNodeType, // 1
			sexpr.ArgumentNodeType, // get_local
			sexpr.ArgumentNodeType, // 2
			sexpr.EndNodeType,      // ) end func
			sexpr.EndNodeType,      // ) end module
		}},
		{
			`(module
				(func (export "add") (param i32 i32) (result i32)
				  local.get 0
				  local.get 1
				  i32.add))`,
			[]sexpr.NodeType{
				sexpr.ModuleNodeType,   // (module
				sexpr.FuncNodeType,     // (func
				sexpr.ExportNodeType,   // (export
				sexpr.ArgumentNodeType, // "add"
				sexpr.EndNodeType,      // ) end export
				sexpr.ParamNodeType,    // (param
				sexpr.ArgumentNodeType, // i32
				sexpr.ArgumentNodeType, // i32
				sexpr.EndNodeType,      // ) end param
				sexpr.ResultNodeType,   // (result
				sexpr.ArgumentNodeType, // i32
				sexpr.EndNodeType,      // ) end result
				sexpr.ArgumentNodeType, // local.get
				sexpr.ArgumentNodeType, // 0
				sexpr.ArgumentNodeType, // local.get
				sexpr.ArgumentNodeType, // 1
				sexpr.ArgumentNodeType, // i32.add
				sexpr.EndNodeType,      // ) end func
				sexpr.EndNodeType,      // ) end module
			},
		},
	}
	for _, test := range tests {
		reader := sexpr.NewReader(strings.NewReader(test.Text))
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
