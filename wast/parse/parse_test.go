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

func TestParse(t *testing.T) {
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
