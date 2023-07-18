// package wit covers parsing and generating wit files
// https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
package wit

import (
	"fmt"
	"io"

	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-wasm/abi/types"
	"github.com/patrickhuber/go-wasm/wit/ast"
	"github.com/patrickhuber/go-wasm/wit/token"
)

func Parse(reader io.Reader) (*ast.Ast, error) {
	lexer := NewLexer(reader)
	return parseAst(lexer)
}

func parseAst(lexer Lexer) (*ast.Ast, error) {
	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}

	n := &ast.Ast{}

	switch tok.Capture {
	case "package":
		packageName, err := parsePackageName(lexer)
		if err != nil {
			return nil, err
		}
		n.PackageName = option.Some(*packageName)
		tok, err = next(lexer)
		if err != nil {
			return nil, err
		}
	default:
		n.PackageName = option.None[ast.PackageName]()
	}

	for {
		item := &ast.AstItem{}
		switch tok.Capture {
		case "use":
			topLevelUse, err := parseTopLevelUse(lexer)
			if err != nil {
				return nil, err
			}
			item.Use = topLevelUse
		case "world":
			world, err := parseWorld(lexer)
			if err != nil {
				return nil, err
			}
			item.World = world
		case "interface":
			inter, err := parseInterface(lexer)
			if err != nil {
				return nil, err
			}
			item.Interface = inter
		}
		n.Items = append(n.Items, *item)
		tok, err = next(lexer)
		if err != nil {
			return nil, err
		}
		if tok.Type == token.EndOfStream {
			break
		}
	}
	return n, nil
}

func parsePackageName(lexer Lexer) (*ast.PackageName, error) {
	packageName := &ast.PackageName{}

	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, err
	}
	packageName.Namespace = tok.Capture

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.Colon)
	if err != nil {
		return nil, err
	}

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, err
	}
	packageName.Name = tok.Capture

	peek, err := peek(lexer)
	if err != nil {
		return nil, err
	}
	if !optional(peek, token.At) {
		packageName.Version = option.None[ast.Version]()
		return packageName, nil
	}

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.At)
	if err != nil {
		return nil, err
	}

	version, err := parseVersion(lexer)
	if err != nil {
		return nil, err
	}
	packageName.Version = option.Some(*version)
	return packageName, nil
}

func parseVersion(lexer Lexer) (*ast.Version, error) {
	version := &ast.Version{}
	return version, nil
}

func parseTopLevelUse(lexer Lexer) (*ast.TopLevelUse, error) {
	return nil, nil
}

func parseInterface(lexer Lexer) (*ast.Interface, error) {
	inter := &ast.Interface{}

	// id
	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, fmt.Errorf("expected id")
	}
	inter.Name = tok.Capture

	// '{'
	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.OpenBrace)
	if err != nil {
		return nil, err
	}

	for {

		interfaceItem, err := parseInterfaceItem(lexer)
		if err != nil {
			return nil, err
		}

		inter.Items = append(inter.Items, *interfaceItem)

		peekTok, err := peek(lexer)
		if err != nil {
			return nil, err
		}

		// exit on '}'
		if optional(peekTok, token.CloseBrace) {
			break
		}
	}

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}

	// '}'
	err = expect(tok, token.CloseBrace)
	if err != nil {
		return nil, err
	}

	return inter, nil
}

func parseInterfaceItem(lexer Lexer) (*ast.InterfaceItem, error) {
	item := &ast.InterfaceItem{}

	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, err
	}
	switch tok.Capture {
	case "use":
	case "resource":
	case "variant":
	case "record":
	case "union":
	case "flags":
	case "enum":
	case "type":

	default:
		// tok == id
		funcItem, err := parseNamedFunc(tok, lexer)
		if err != nil {
			return nil, err
		}
		item.Func = funcItem
	}
	return item, nil
}

func parseNamedFunc(name *token.Token, lexer Lexer) (*ast.NamedFunc, error) {
	named := &ast.NamedFunc{
		Name: name.Capture,
	}

	tok, err := lexer.Next()
	if err != nil {
		return nil, err
	}

	err = expect(tok, token.Colon)
	if err != nil {
		return nil, err
	}

	_func, err := parseFunc(lexer)
	if err != nil {
		return nil, err
	}
	named.Func = _func
	return named, nil
}

func parseFunc(lexer Lexer) (*ast.Func, error) {
	_func := &ast.Func{}

	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}

	err = expectValue(tok, token.String, "func")
	if err != nil {
		return nil, err
	}

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}

	err = expect(tok, token.OpenParen)
	if err != nil {
		return nil, err
	}

	for {
		// name ':' type
		param, err := parseParameter(lexer)
		if err != nil {
			return nil, err
		}

		_func.Params = append(_func.Params, *param)

		peekTok, err := peek(lexer)
		if err != nil {
			return nil, err
		}

		if peekTok.Type == token.Comma {
			tok, err = next(lexer)
			if err != nil {
				return nil, err
			}
			err = expect(tok, token.Comma)
			if err != nil {
				return nil, err
			}
		} else if peekTok.Type == token.CloseParen {
			break
		} else {
			return nil, fmt.Errorf("%w. expected ',' or ')' but found %s", parseError(peekTok), peekTok.Capture)
		}

		peekTok, err = peek(lexer)
		if err != nil {
			return nil, err
		}
		if peekTok.Type == token.CloseParen {
			break
		}
	}

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.CloseParen)
	if err != nil {
		return nil, err
	}

	return _func, nil
}

func parseParameter(lexer Lexer) (*ast.Parameter, error) {
	parameter := &ast.Parameter{}

	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, err
	}
	parameter.Id = tok.Capture

	tok, err = next(lexer)
	if err != nil {
		return nil, err
	}

	err = expect(tok, token.Colon)
	if err != nil {
		return nil, err
	}

	ty, err := parseType(lexer)
	if err != nil {
		return nil, err
	}
	parameter.Type = ty
	return parameter, nil
}

func parseType(lexer Lexer) (types.Type, error) {
	tok, err := next(lexer)
	if err != nil {
		return nil, err
	}
	err = expect(tok, token.String)
	if err != nil {
		return nil, err
	}
	var ty types.Type
	switch tok.Capture {
	case "string":
		ty = types.NewString()
	default:
		return nil, fmt.Errorf("unrecognized type %s", tok.Capture)
	}
	return ty, nil
}

func parseWorld(lexer Lexer) (*ast.World, error) {
	return nil, nil
}

func next(lexer Lexer) (*token.Token, error) {
	for {
		tok, err := lexer.Next()
		if err != nil {
			return nil, err
		}
		if tok.Type != token.Whitespace {
			return tok, nil
		}
	}
}

func peek(lexer Lexer) (*token.Token, error) {
	for {
		peek, err := lexer.Peek()
		if err != nil {
			return nil, err
		}
		if peek.Type != token.Whitespace {
			return peek, nil
		}
		_, err = lexer.Next()
		if err != nil {
			return nil, err
		}
	}
}

func optional(tok *token.Token, tokenType token.TokenType) bool {
	return tok.Type == tokenType
}

func optionalValue(tok *token.Token, tokenType token.TokenType, capture string) bool {
	return optional(tok, tokenType) && tok.Capture == capture
}

func expect(tok *token.Token, tokenType token.TokenType) error {
	if optional(tok, tokenType) {
		return nil
	}
	return fmt.Errorf("%w. expected %v but found %v", parseError(tok), tokenType, tok.Type)
}

func expectValue(tok *token.Token, tokenType token.TokenType, capture string) error {
	if optionalValue(tok, tokenType, capture) {
		return nil
	}
	return fmt.Errorf("%w. expected %v but found %v", parseError(tok), tokenType, tok.Type)
}

func parseError(tok *token.Token) error {
	return fmt.Errorf("error parsing at position %d line %d column %d", tok.Position, tok.Line, tok.Column)
}
