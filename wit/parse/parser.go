// package wit covers parsing and generating wit files
// https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
package wit

import (
	"fmt"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wit/ast"
	"github.com/patrickhuber/go-wasm/wit/lex"
	"github.com/patrickhuber/go-wasm/wit/token"
)

func Parse(input string) (*ast.Ast, error) {
	lexer := lex.New(input)
	return parseAst(lexer).Deconstruct()
}

func parseAst(lexer *lex.Lexer) (res types.Result[*ast.Ast]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	n := &ast.Ast{}

	switch tok.Type {
	case token.Package:
		packageName := parsePackageName(lexer).Unwrap()
		n.PackageName = option.Some(*packageName)
		tok = next(lexer).Unwrap()
	default:
		n.PackageName = option.None[ast.PackageName]()
	}

	for tok.Type != token.EndOfStream {
		item := &ast.AstItem{}
		switch tok.Type {
		case token.Use:
			item.Use = parseTopLevelUse(lexer).Unwrap()
		case token.World:
			item.World = parseWorld(lexer).Unwrap()
		case token.Interface:
			item.Interface = parseInterface(lexer).Unwrap()
		}
		n.Items = append(n.Items, *item)

		tok = next(lexer).Unwrap()
	}
	return result.Ok(n)
}

func parsePackageName(lexer *lex.Lexer) (res types.Result[*ast.PackageName]) {
	defer handle.Error(&res)

	// id
	packageName := &ast.PackageName{}
	packageName.Name = next(lexer).Unwrap().Capture

	// ':'
	expect(lexer, token.Colon).Unwrap()

	// id
	packageName.Name = parseId(lexer).Unwrap()

	// @
	if !eat(lexer, token.At).Unwrap() {
		packageName.Version = option.None[ast.Version]()
		return result.Ok(packageName)
	}

	version := parseVersion(lexer).Unwrap()
	packageName.Version = option.Some(*version)

	return result.Ok(packageName)
}

func parseVersion(lexer *lex.Lexer) (res types.Result[*ast.Version]) {
	defer handle.Error(&res)
	version := &ast.Version{}
	return result.Ok(version)
}

func parseTopLevelUse(lexer *lex.Lexer) (res types.Result[*ast.TopLevelUse]) {
	defer handle.Error(&res)

	topLevelUse := &ast.TopLevelUse{
		Item: parseUsePath(lexer).Unwrap(),
	}
	if eat(lexer, token.As).Unwrap() {
		topLevelUse.As = option.Some(parseId(lexer).Unwrap())
	} else {
		topLevelUse.As = option.None[string]()
	}
	return result.Ok(topLevelUse)
}

func parseInterface(lexer *lex.Lexer) (res types.Result[*ast.Interface]) {

	defer handle.Error(&res)
	inter := &ast.Interface{}

	// id
	inter.Name = parseId(lexer).Unwrap()

	// '{'
	expect(lexer, token.OpenBrace).Unwrap()

	// exit on '}'
	inter.Items = parseInterfaceItems(lexer).Unwrap()
	return result.Ok(inter)
}

func parseInterfaceItems(lexer *lex.Lexer) (res types.Result[[]ast.InterfaceItem]) {
	defer handle.Error(&res)
	var items []ast.InterfaceItem
	for !eat(lexer, token.CloseBrace).Unwrap() {
		items = append(items, *parseInterfaceItem(lexer).Unwrap())
	}
	return result.Ok(items)
}

func parseInterfaceItem(lexer *lex.Lexer) (res types.Result[*ast.InterfaceItem]) {
	defer handle.Error(&res)

	itemType := next(lexer).Unwrap()
	item := &ast.InterfaceItem{}

	switch itemType.Type {
	case token.Use:
		item.Use = parseUse(lexer).Unwrap()
	case token.Resource:
		item.TypeDef = parseResource(lexer).Unwrap()
	case token.Record:
		item.TypeDef = parseRecord(lexer).Unwrap()
	case token.Flags:
		item.TypeDef = parseFlags(lexer).Unwrap()
	case token.Variant:
		item.TypeDef = parseVariant(lexer).Unwrap()
	case token.Union:
		item.TypeDef = parseUnion(lexer).Unwrap()
	case token.Enum:
		item.TypeDef = parseEnum(lexer).Unwrap()
	case token.Type:
		item.TypeDef = parseTypeDef(lexer).Unwrap()

	default:
		// tok == id
		item.Func = parseNamedFunc(itemType, lexer).Unwrap()
	}
	return result.Ok(item)
}

func parseTypeDef(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)

	name := parseId(lexer).Unwrap()
	expect(lexer, token.Equal).Unwrap()
	ty := parseType(lexer).Unwrap()
	return result.Ok(
		&ast.TypeDef{
			Name: name,
			Type: ty,
		})
}

func parseResource(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)

	name := parseId(lexer).Unwrap()
	resource := &ast.Resource{}
	if eat(lexer, token.OpenBrace).Unwrap() {
		for {
			if eat(lexer, token.CloseBrace).Unwrap() {
				break
			}
			resourceFunc := parseResourceFunc(lexer).Unwrap()
			resource.Functions = append(resource.Functions, resourceFunc)
		}
	}

	return result.Ok(&ast.TypeDef{
		Name: name,
		Type: resource,
	})
}

/*
resource-method ::= func-item
| id ':' 'static' func-type
| id ':' func-type
| 'constructor' param-list
*/
func parseResourceFunc(lexer *lex.Lexer) (res types.Result[ast.ResourceFunc]) {
	defer handle.Error(&res)

	namedFunc := &ast.NamedFunc{
		Func: &ast.Func{},
	}
	if eat(lexer, token.Constructor).Unwrap() {
		expect(lexer, token.OpenParen).Unwrap()
		namedFunc.Func.Params = parseParameters(lexer).Unwrap()
		namedFunc.Name = "constructor"
		return result.Ok[ast.ResourceFunc](&ast.Constructor{
			NamedFunc: namedFunc,
		})
	}

	namedFunc.Name = parseId(lexer).Unwrap()
	expect(lexer, token.Colon).Unwrap()

	static := eat(lexer, token.Static).Unwrap()

	expect(lexer, token.Func).Unwrap()
	namedFunc.Func = parseFunc(lexer).Unwrap()

	if static {
		return result.Ok[ast.ResourceFunc](&ast.Static{
			NamedFunc: namedFunc,
		})
	}
	return result.Ok[ast.ResourceFunc](&ast.Method{
		NamedFunc: namedFunc,
	})
}

func parseUse(lexer *lex.Lexer) (res types.Result[*ast.Use]) {
	defer handle.Error(&res)
	u := &ast.Use{
		From: parseUsePath(lexer).Unwrap(),
	}

	// .
	expect(lexer, token.Period).Unwrap()
	u.Names = parseItemList[ast.UseName](
		lexer,
		token.OpenBrace,
		token.CloseBrace,
		parseUseName).Unwrap()

	return result.Ok(u)
}

func parseUseName(lexer *lex.Lexer) types.Result[ast.UseName] {
	name := ast.UseName{
		Name: parseId(lexer).Unwrap(),
		As:   option.None[string](),
	}

	// as
	if eat(lexer, token.As).Unwrap() {
		name.As = option.Some(parseId(lexer).Unwrap())
	}
	return result.Ok(name)
}

func parseItemList[T any](
	lexer *lex.Lexer,
	begin token.TokenType,
	end token.TokenType,
	parseItem func(l *lex.Lexer) types.Result[T]) (res types.Result[[]T]) {

	defer handle.Error(&res)

	var itemList []T
	expect(lexer, begin).Unwrap()
	for !eat(lexer, end).Unwrap() {

		item := parseItem(lexer).Unwrap()
		itemList = append(itemList, item)

		if !eat(lexer, token.Comma).Unwrap() {
			expect(lexer, end).Unwrap()
			break
		}
	}
	return result.Ok(itemList)
}

func parseUsePath(lexer *lex.Lexer) (res types.Result[*ast.UsePath]) {
	defer handle.Error(&res)
	id := parseId(lexer).Unwrap()

	// `foo`
	if !eat(lexer, token.Colon).Unwrap() {
		return result.Ok(&ast.UsePath{
			Id: id,
		})
	}

	// `foo:bar/baz@1.0`
	return parsePath(lexer, id)
}

func parsePath(lexer *lex.Lexer, namespace string) (res types.Result[*ast.UsePath]) {
	defer handle.Error(&res)

	pkgName := parseId(lexer).Unwrap()
	expect(lexer, token.Slash).Unwrap()
	name := parseId(lexer).Unwrap()
	version := parseOptionalVersion(lexer).Unwrap()
	return result.Ok(&ast.UsePath{
		Package: struct {
			Id   *ast.PackageName
			Name string
		}{
			Id: &ast.PackageName{
				Namespace: namespace,
				Name:      pkgName,
				Version:   version,
			},
			Name: name,
		},
	})
}

func parseOptionalVersion(lexer *lex.Lexer) (res types.Result[types.Option[ast.Version]]) {
	return result.Errorf[types.Option[ast.Version]]("not implemented")
}

func parseNamedFunc(name *token.Token, lexer *lex.Lexer) (res types.Result[*ast.NamedFunc]) {

	defer handle.Error(&res)

	named := &ast.NamedFunc{
		Name: name.Capture,
	}

	expect(lexer, token.Colon).Unwrap()
	expect(lexer, token.Func).Unwrap()

	_func := parseFunc(lexer).Unwrap()
	named.Func = _func

	return result.Ok(named)
}

func parseFunc(lexer *lex.Lexer) (res types.Result[*ast.Func]) {

	defer handle.Error(&res)

	// (
	expect(lexer, token.OpenParen).Unwrap()

	parameters := parseParameters(lexer).Unwrap()
	results := &ast.ResultList{}
	if eat(lexer, token.RightArrow).Unwrap() {
		if eat(lexer, token.OpenParen).Unwrap() {
			results.Named = parseParameters(lexer).Unwrap()
		} else {
			results.Anonymous = parseType(lexer).Unwrap()
		}
	} else {
		results.Named = nil // ? []ast.Parameter{}
	}

	return result.Ok(&ast.Func{
		Params:  parameters,
		Results: results,
	})
}

func parseParameters(lexer *lex.Lexer) (res types.Result[[]ast.Parameter]) {
	var parameters []ast.Parameter
	for {

		// )
		if eat(lexer, token.CloseParen).Unwrap() {
			break
		}

		// name ':' type
		param := parseParameter(lexer).Unwrap()
		parameters = append(parameters, *param)

		peekTok := peek(lexer).Unwrap()

		if peekTok.Type == token.Comma {
			expect(lexer, token.Comma).Unwrap()
		} else if peekTok.Type == token.CloseParen {
			expect(lexer, token.CloseParen).Unwrap()
			break
		} else {
			return result.Errorf[[]ast.Parameter]("%w. expected ',' or ')' but found %s", parseError(peekTok), peekTok.Capture)
		}
	}
	return result.Ok(parameters)
}

func parseParameter(lexer *lex.Lexer) (res types.Result[*ast.Parameter]) {
	defer handle.Error(&res)

	parameter := &ast.Parameter{}
	parameter.Id = parseId(lexer).Unwrap()

	expect(lexer, token.Colon).Unwrap()

	parameter.Type = parseType(lexer).Unwrap()

	return result.Ok(parameter)
}

func parseType(lexer *lex.Lexer) (res types.Result[ast.Type]) {
	defer handle.Error(&res)

	name := next(lexer).Unwrap()

	var ty ast.Type
	switch name.Type {
	case token.U32:
		ty = &ast.U32{}
	case token.String:
		ty = &ast.String{}
	case token.Float32:
		ty = &ast.Float32{}
	case token.Float64:
		ty = &ast.Float64{}
	case token.Stream:
		ty = parseStream(lexer).Unwrap()
	case token.Future:
		ty = parseFuture(lexer).Unwrap()
	case token.List:
		ty = parseList(lexer).Unwrap()
	case token.Option:
		ty = parseOption(lexer).Unwrap()
	case token.Result:
		ty = parseResult(lexer).Unwrap()
	case token.Tuple:
		ty = parseTuple(lexer).Unwrap()
	case token.Own:
		expect(lexer, token.Less).Unwrap()
		ty = &ast.Own{
			Id: parseId(lexer).Unwrap(),
		}
		expect(lexer, token.Greater).Unwrap()
	case token.Borrow:
		expect(lexer, token.Less).Unwrap()
		ty = &ast.Borrow{
			Id: parseId(lexer).Unwrap(),
		}
		expect(lexer, token.Greater).Unwrap()
	default:
		ty = &ast.Id{Value: name.Capture}
	}

	return result.Ok(ty)
}

// stream<T, Z>
// stream<_, Z>
// stream<T>
// stream
func parseStream(lexer *lex.Lexer) (res types.Result[*ast.Stream]) {
	defer handle.Error(&res)
	stream := &ast.Stream{
		End:     option.None[ast.Type](),
		Element: option.None[ast.Type](),
	}
	if eat(lexer, token.Less).Unwrap() {
		if eat(lexer, token.Underscore).Unwrap() {
			expect(lexer, token.Comma).Unwrap()
			stream.End = option.Some(parseType(lexer).Unwrap())
		} else {
			stream.Element = option.Some(parseType(lexer).Unwrap())
			if eat(lexer, token.Comma).Unwrap() {
				stream.End = option.Some(parseType(lexer).Unwrap())
			}
		}
		expect(lexer, token.Greater).Unwrap()
	}
	return result.Ok(stream)
}

func parseFuture(lexer *lex.Lexer) (res types.Result[*ast.Future]) {
	defer handle.Error(&res)
	future := &ast.Future{}
	if eat(lexer, token.Less).Unwrap() {
		ty := parseType(lexer).Unwrap()
		future.Type = option.Some(ty)
		expect(lexer, token.Greater).Unwrap()
	}
	return result.Ok(future)
}

func parseList(lexer *lex.Lexer) (res types.Result[*ast.List]) {
	defer handle.Error(&res)
	expect(lexer, token.Less).Unwrap()
	ty := parseType(lexer).Unwrap()
	expect(lexer, token.Greater).Unwrap()
	return result.Ok(&ast.List{
		Type: ty,
	})
}

func parseOption(lexer *lex.Lexer) (res types.Result[*ast.Option]) {
	defer handle.Error(&res)
	expect(lexer, token.Less).Unwrap()
	ty := parseType(lexer).Unwrap()
	expect(lexer, token.Greater).Unwrap()
	return result.Ok(&ast.Option{
		Type: ty,
	})
}

func parseTuple(lexer *lex.Lexer) (res types.Result[*ast.Tuple]) {
	defer handle.Error(&res)
	tuple := &ast.Tuple{
		Types: parseItemList[ast.Type](lexer, token.Less, token.Greater, parseType).Unwrap(),
	}
	return result.Ok(tuple)
}

// result<T, E>
// result<_, E>
// result<T>
// result
func parseResult(lexer *lex.Lexer) (res types.Result[*ast.Result]) {
	defer handle.Error(&res)
	r := &ast.Result{
		Ok:    option.None[ast.Type](),
		Error: option.None[ast.Type](),
	}
	if eat(lexer, token.Less).Unwrap() {
		if eat(lexer, token.Underscore).Unwrap() {
			expect(lexer, token.Comma).Unwrap()
			r.Error = option.Some(parseType(lexer).Unwrap())
		} else {
			r.Ok = option.Some(parseType(lexer).Unwrap())
			if eat(lexer, token.Comma).Unwrap() {
				r.Error = option.Some(parseType(lexer).Unwrap())
			}
		}
		expect(lexer, token.Greater).Unwrap()
	}
	return result.Ok(r)
}

func parseWorld(lexer *lex.Lexer) (res types.Result[*ast.World]) {
	defer handle.Error(&res)

	id := parseId(lexer).Unwrap()

	expect(lexer, token.OpenBrace).Unwrap()

	worldItems := parseWorldItems(lexer).Unwrap()
	world := &ast.World{
		Id:    id,
		Items: worldItems,
	}
	return result.Ok(world)
}

func parseWorldItems(lexer *lex.Lexer) (res types.Result[[]ast.WorldItem]) {
	defer handle.Error(&res)
	var worldItems []ast.WorldItem
	for !eat(lexer, token.CloseBrace).Unwrap() {
		worldItems = append(worldItems, parseWorldItem(lexer).Unwrap())
	}
	return result.Ok(worldItems)
}

func parseWorldItem(lexer *lex.Lexer) (res types.Result[ast.WorldItem]) {
	defer handle.Error(&res)

	itemType := next(lexer).Unwrap()
	switch itemType.Type {
	case token.Export:
		return parseExport(lexer)
	case token.Import:
		return parseImport(lexer)
	case token.Use:
		return result.Ok[ast.WorldItem](
			parseUse(lexer).Unwrap(),
		)
	case token.Type:
		return result.Ok[ast.WorldItem](
			parseTypeDef(lexer).Unwrap(),
		)
	case token.Record:
		return result.Ok[ast.WorldItem](
			parseRecord(lexer).Unwrap(),
		)
	case token.Variant:
		return result.Ok[ast.WorldItem](
			parseVariant(lexer).Unwrap(),
		)
	case token.Resource:
		return result.Ok[ast.WorldItem](
			parseResource(lexer).Unwrap(),
		)
	case token.Include:
		return parseInclude(lexer)
	}
	return result.Errorf[ast.WorldItem]("%w : Unrecognized world item '%s'. Expected (export, import, resource, use, type, include). Found token.%v",
		parseErrorFromLexer(lexer),
		itemType.Capture,
		itemType.Type)
}

func parseExport(lexer *lex.Lexer) (res types.Result[ast.WorldItem]) {
	defer handle.Error(&res)

	// this ID can have different meanings depending on what follows
	id := parseId(lexer).Unwrap()

	return result.Ok[ast.WorldItem](&ast.ExportExternType{
		ExternType: parseExternType(lexer, id).Unwrap(),
	})
}

func parseImport(lexer *lex.Lexer) (res types.Result[ast.WorldItem]) {
	defer handle.Error(&res)

	// this ID can have different meanings depending on what follows
	id := parseId(lexer).Unwrap()

	return result.Ok[ast.WorldItem](&ast.ImportExternType{
		ExternType: parseExternType(lexer, id).Unwrap(),
	})
}

func parseInclude(lexer *lex.Lexer) (res types.Result[ast.WorldItem]) {
	defer handle.Error(&res)
	include := &ast.Include{
		From: parseUsePath(lexer).Unwrap(),
	}
	if eat(lexer, token.With).Unwrap() {
		expect(lexer, token.OpenBrace).Unwrap()
		for !eat(lexer, token.CloseBrace).Unwrap() {
			id := parseId(lexer).Unwrap()
			expect(lexer, token.As).Unwrap()
			as := parseId(lexer).Unwrap()
			include.Names = append(include.Names, ast.IncludeName{
				Name: id,
				As:   as,
			})
		}
	}
	return result.Ok[ast.WorldItem](include)
}

func parseExternType(lexer *lex.Lexer, id string) (res types.Result[*ast.ExternType]) {
	defer handle.Error(&res)

	et := &ast.ExternType{}
	if eat(lexer, token.Colon).Unwrap() {
		if eat(lexer, token.Func).Unwrap() {
			// import foo: func(...)
			//                 ^
			et.Func = parseFunc(lexer).Unwrap()
		} else if eat(lexer, token.Interface).Unwrap() {
			// import foo: interface{...}
			//                      ^
			expect(lexer, token.OpenBrace).Unwrap()
			et.Interface = &ast.Interface{
				Items: parseInterfaceItems(lexer).Unwrap(),
			}
		}
	} else {
		// import foo
		//           ^
		et.UsePath = &ast.UsePath{
			Id: id,
		}
	}
	return result.Ok(et)
}

// record ::= id '{' fields '}'
// fields ::= field | field ',' fields | λ
// field  ::= id ':' ty
func parseRecord(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)

	name := parseId(lexer).Unwrap()
	record := &ast.Record{
		Fields: parseItemList(lexer, token.OpenBrace, token.CloseBrace,
			func(l *lex.Lexer) types.Result[ast.Field] {
				id := parseId(lexer).Unwrap()
				expect(lexer, token.Colon).Unwrap()
				ty := parseType(lexer).Unwrap()
				return result.Ok(ast.Field{Name: id, Type: ty})
			}).Unwrap(),
	}

	return result.Ok(&ast.TypeDef{
		Name: name,
		Type: record,
	})
}

func parseFlags(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)

	name := parseId(lexer).Unwrap()
	flags := &ast.Flags{
		Flags: parseItemList(lexer, token.OpenBrace, token.CloseBrace, func(l *lex.Lexer) types.Result[ast.Flag] {
			id := parseId(lexer).Unwrap()
			return result.Ok(ast.Flag{
				Id: id,
			})
		}).Unwrap(),
	}
	return result.Ok(&ast.TypeDef{
		Name: name,
		Type: flags,
	})
}

func parseVariant(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)

	name := parseId(lexer).Unwrap()

	cases := parseItemList(lexer, token.OpenBrace, token.CloseBrace, func(l *lex.Lexer) types.Result[ast.Case] {
		name := parseId(lexer).Unwrap()
		c := &ast.Case{
			Name: name,
			Type: option.None[ast.Type](),
		}
		if eat(lexer, token.OpenParen).Unwrap() {
			ty := parseType(lexer).Unwrap()
			expect(lexer, token.CloseParen)
			c.Type = option.Some(ty)
		}
		return result.Ok(*c)
	}).Unwrap()
	return result.Ok(&ast.TypeDef{
		Name: name,
		Type: &ast.Variant{
			Cases: cases,
		},
	})
}

func parseUnion(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)
	id := parseId(lexer).Unwrap()
	cases := parseItemList(lexer, token.OpenBrace, token.CloseBrace, func(l *lex.Lexer) types.Result[ast.UnionCase] {
		ty := parseType(lexer).Unwrap()
		return result.Ok(ast.UnionCase{
			Type: ty,
		})
	}).Unwrap()
	union := &ast.Union{
		Cases: cases,
	}
	return result.Ok(&ast.TypeDef{
		Name: id,
		Type: union,
	})
}

func parseEnum(lexer *lex.Lexer) (res types.Result[*ast.TypeDef]) {
	defer handle.Error(&res)
	id := parseId(lexer).Unwrap()
	cases := parseItemList(lexer, token.OpenBrace, token.CloseBrace,
		func(l *lex.Lexer) types.Result[ast.EnumCase] {
			id := parseId(lexer).Unwrap()
			return result.Ok(ast.EnumCase{
				Name: id,
			})
		}).Unwrap()
	enum := &ast.Enum{
		Cases: cases,
	}
	return result.Ok(&ast.TypeDef{
		Name: id,
		Type: enum,
	})
}

func parseId(lexer *lex.Lexer) (res types.Result[string]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	switch tok.Type {
	case token.Id:
		return result.Ok(tok.Capture)
	case token.ExplicitId:
		return result.Ok(tok.Capture)
	default:
		return result.Errorf[string]("%w : found value '%s', type '%v' but expected (id, explicit_id)", parseError(tok), tok.Capture, tok.Type)
	}
}

func eat(lexer *lex.Lexer, tokenType token.TokenType) (res types.Result[bool]) {
	defer handle.Error(&res)

	tok := peek(lexer).Unwrap()
	if !is(tok, tokenType) {
		return result.Ok(false)
	}

	expect(lexer, tokenType).Unwrap()
	return result.Ok(true)
}

func expect(lexer *lex.Lexer, tokenType token.TokenType) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type == tokenType {
		return result.Ok[any](nil)
	}
	return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), tokenType, tok.Type)
}

func peek(lexer *lex.Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		p := result.New(lexer.Peek())
		r := p.Unwrap()
		if r.Type != token.Whitespace {
			return p
		}
		// consume whitespace
		_ = result.New(lexer.Next()).Unwrap()
	}
}

func next(lexer *lex.Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		res = result.New(lexer.Next())
		tok := res.Unwrap()
		switch tok.Type {
		// skip whitespace
		case token.Whitespace:
			continue
		// skip comments
		case token.BlockComment:
			continue
		// skip comments
		case token.LineComment:
			continue
		}
		return
	}
}

func is(tok *token.Token, tokenType token.TokenType) bool {
	return tok.Type == tokenType
}

func parseError(tok *token.Token) error {
	line := tok.Line + 1
	col := tok.Column + 1
	return fmt.Errorf(
		"error parsing at line %d, column %d, position %d",
		line,
		col,
		tok.Position)
}

func parseErrorFromLexer(lexer *lex.Lexer) error {
	line := lexer.Line() + 1
	col := lexer.Column() + 1
	return fmt.Errorf(
		"error parsing at line %d, column %d",
		line,
		col)
}
