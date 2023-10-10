package parse_test

import (
	"bufio"
	"io"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-wasm/wast/ast"
	"github.com/patrickhuber/go-wasm/wast/parse"
	wat "github.com/patrickhuber/go-wasm/wat/ast"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	dir := "../../submodules/github.com/WebAssembly/testsuite"
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.NotNil(t, files)

	filter, err := regexp.Compile("i32.wast")
	require.NoError(t, err)

	for _, file := range files {
		if !filter.MatchString(file.Name()) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {

			full := path.Join(dir, file.Name())
			file, err := os.OpenFile(full, os.O_RDONLY, 0666)
			require.NoError(t, err)

			reader := bufio.NewReader(file)
			bytes, err := io.ReadAll(reader)
			require.NoError(t, err)

			input := string(bytes)
			directives, err := parse.Parse(input)
			require.NoError(t, err)
			require.Greater(t, 0, len(directives))
		})
	}
}

func TestParse(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected []ast.Directive
	}
	tests := []test{
		{"assert_return", `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`, []ast.Directive{ast.AssertReturn{
			Action: ast.Invoke{
				String: "add",
				Name:   option.None[string](),
				Const: []ast.Const{
					ast.I32Const{Value: 1},
					ast.I32Const{Value: 1},
				},
			},
			Results: []ast.Result{
				ast.I32Const{Value: 2},
			},
		}}},
		{"assert_trap", `(assert_trap (invoke "div_s" (i32.const 1) (i32.const 0)) "integer divide by zero")`, []ast.Directive{ast.AssertTrap{
			Action: ast.Invoke{
				String: "div_s",
				Name:   option.None[string](),
				Const: []ast.Const{
					ast.I32Const{Value: 1},
					ast.I32Const{Value: 0},
				},
			},
			Failure: "integer divide by zero",
		}}},
		{"assert_invalid", `(assert_invalid
			(module
			  (func $type-unary-operand-empty
				(i32.eqz) (drop)
			  )
			)
			"type mismatch"
		  )`, []ast.Directive{
			ast.AssertInvalid{
				Module: &wat.Module{
					Functions: []wat.Function{
						{
							ID: option.Some("$type-unary-operand-empty"),
							Instructions: []wat.Instruction{
								wat.I32Eqz{},
								wat.Drop{},
							},
						},
					},
				},
				Failure: "type mismatch",
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directives, err := parse.Parse(test.input)
			require.NoError(t, err)
			require.Equal(t, test.expected, directives)
		})
	}
}
