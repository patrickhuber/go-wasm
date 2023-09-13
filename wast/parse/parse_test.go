package parse_test

import (
	"bufio"
	"io"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/patrickhuber/go-wasm/wast/parse"

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
		name  string
		input string
	}
	tests := []test{
		{"assert_return", `(assert_return (invoke "add" (i32.const 1) (i32.const 1)) (i32.const 2))`},
		{"assert_trap", `(assert_trap (invoke "div_s" (i32.const 1) (i32.const 0)) "integer divide by zero")`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directives, err := parse.Parse(test.input)
			require.NoError(t, err)
			require.Greater(t, 0, len(directives))
		})
	}
}
