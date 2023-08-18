package wit_test

import (
	"os"
	"path"
	"regexp"
	"testing"

	wit "github.com/patrickhuber/go-wasm/wit/parse"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	dir := "../../submodules/github.com/bytecodealliance/wasm-tools/crates/wit-parser/tests/ui"
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.NotNil(t, files)

	witFilter, err := regexp.Compile(".*[.]wit$")
	require.NoError(t, err)

	for _, file := range files {

		if !witFilter.MatchString(file.Name()) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			full := path.Join(dir, file.Name())
			bytes, err := os.ReadFile(full)
			require.NoError(t, err)

			node, err := wit.Parse(string(bytes))
			require.NoError(t, err)
			require.NotNil(t, node)
		})
	}
}

func TestParseFail(t *testing.T) {
	dir := "../../submodules/github.com/bytecodealliance/wasm-tools/crates/wit-parser/tests/ui/parse-fail"
	files, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.NotNil(t, files)

	witFilter, err := regexp.Compile(".*[.]wit$")
	require.NoError(t, err)

	for _, file := range files {

		if !witFilter.MatchString(file.Name()) {
			continue
		}
		t.Run(file.Name(), func(t *testing.T) {
			full := path.Join(dir, file.Name())
			bytes, err := os.ReadFile(full)
			require.NoError(t, err)
			_, err = wit.Parse(string(bytes))
			require.Error(t, err)
		})
	}
}
