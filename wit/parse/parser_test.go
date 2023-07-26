package wit_test

import (
	"bufio"
	"errors"
	"io"
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

			file, err := os.OpenFile(full, os.O_RDONLY, 0666)
			require.NoError(t, err)
			reader := bufio.NewReader(file)

			var runes []rune
			for {
				r, _, err := reader.ReadRune()
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				runes = append(runes, r)
			}

			node, err := wit.Parse(runes)
			require.NoError(t, err)
			require.NotNil(t, node)
		})
	}

}
