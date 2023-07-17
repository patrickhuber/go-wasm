package wit_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/go-wasm/wit"
	"github.com/patrickhuber/go-wasm/wit/token"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	t.Run("line_comment", func(t *testing.T) {
		str := "// this is a comment line"
		lex := wit.NewLexer(strings.NewReader(str))
		tok, err := lex.Next()
		require.Nil(t, err)
		require.NotNil(t, tok)
		require.Equal(t, str, tok.Capture)
		require.Equal(t, token.LineComment, tok.Type)
	})
	t.Run("block_comment", func(t *testing.T) {
		str := "/* this is a comment block */"
		lex := wit.NewLexer(strings.NewReader(str))
		tok, err := lex.Next()
		require.Nil(t, err)
		require.NotNil(t, tok)
		require.Equal(t, str, tok.Capture)
		require.Equal(t, token.BlockComment, tok.Type)
	})
	t.Run("whitespace", func(t *testing.T) {
		str := "\f\t "
		lex := wit.NewLexer(strings.NewReader(str))
		tok, err := lex.Next()
		require.Nil(t, err)
		require.NotNil(t, tok)
		require.Equal(t, str, tok.Capture)
		require.Equal(t, token.Whitespace, tok.Type)
	})
}
