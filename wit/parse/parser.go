// package wit covers parsing and generating wit files
// https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md
package wit

import (
	"fmt"
	"strconv"

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
		packageDeclaration := parsePackageDeclaration(lexer).Unwrap()
		n.PackageDeclaration = option.Some(*packageDeclaration)
		tok = next(lexer).Unwrap()
	default:
		n.PackageDeclaration = option.None[ast.PackageDeclaration]()
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

func parsePackageDeclaration(lexer *lex.Lexer) (res types.Result[*ast.PackageDeclaration]) {
	defer handle.Error(&res)

	// id
	packageName := &ast.PackageDeclaration{}
	packageName.Namespace = next(lexer).Unwrap().Capture

	// ':'
	expect(lexer, token.Colon).Unwrap()

	// id
	packageName.Name = parseId(lexer).Unwrap()

	// @
	if !eat(lexer, token.At).Unwrap() {
		packageName.Version = option.None[ast.Version]()
	} else {
		version := parseVersion(lexer).Unwrap()
		packageName.Version = option.Some(*version)
	}
	// ;
	expect(lexer, token.Semicolon).Unwrap()
	return result.Ok(packageName)
}

func parseVersion(lexer *lex.Lexer) (res types.Result[*ast.Version]) {
	defer handle.Error(&res)

	major := parseInteger(lexer).Unwrap()
	expect(lexer, token.Period).Unwrap()

	minor := parseInteger(lexer).Unwrap()
	expect(lexer, token.Period).Unwrap()

	patch := parseInteger(lexer).Unwrap()

	return result.Ok(&ast.Version{
		Major: uint64(major),
		Minor: uint64(minor),
		Patch: uint64(patch),
	})
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

	// '{' interface-items '}'
	inter.Items = parseInterfaceItems(lexer).Unwrap()

	return result.Ok(inter)
}

func parseInterfaceItems(lexer *lex.Lexer) (res types.Result[[]ast.InterfaceItem]) {
	defer handle.Error(&res)
	var items []ast.InterfaceItem
	expect(lexer, token.OpenBrace).Unwrap()
	for !eat(lexer, token.CloseBrace).Unwrap() {
		items = append(items, parseInterfaceItem(lexer).Unwrap())
	}
	return result.Ok(items)
}

func parseInterfaceItem(lexer *lex.Lexer) (res types.Result[ast.InterfaceItem]) {
	defer handle.Error(&res)

	itemType := peek(lexer).Unwrap()
	var item ast.InterfaceItem

	switch itemType.Type {
	case token.Use:
		item = parseUse(lexer).Unwrap()
	case token.Resource:
		item = parseResource(lexer).Unwrap()
	case token.Record:
		item = parseRecord(lexer).Unwrap()
	case token.Flags:
		item = parseFlags(lexer).Unwrap()
	case token.Variant:
		item = parseVariant(lexer).Unwrap()
	case token.Enum:
		item = parseEnum(lexer).Unwrap()
	case token.Type:
		item = parseTypeDef(lexer).Unwrap()

	default:
		// tok == id
		item = parseFuncItem(lexer).Unwrap()
	}
	return result.Ok(item)
}

func parseTypeDef(lexer *lex.Lexer) (res types.Result[ast.TypeDef]) {
	defer handle.Error(&res)

	ty := peek(lexer).Unwrap()
	var typeDef ast.TypeDef

	// 'resource' | 'variant' | 'record' | 'flags' | 'enum' | 'type'
	switch ty.Type {
	case token.Resource:
		typeDef = parseResource(lexer).Unwrap()
	case token.Variant:
		typeDef = parseVariant(lexer).Unwrap()
	case token.Record:
		typeDef = parseRecord(lexer).Unwrap()
	case token.Flags:
		typeDef = parseFlags(lexer).Unwrap()
	case token.Enum:
		typeDef = parseEnum(lexer).Unwrap()
	case token.Type:
		typeDef = parseTypeItem(lexer).Unwrap()
	default:
		return result.Errorf[ast.TypeDef]("error parsing TypeDef %w", parseError(ty))
	}
	return result.Ok(typeDef)
}

func parseResource(lexer *lex.Lexer) (res types.Result[ast.Resource]) {
	defer handle.Error(&res)

	expect(lexer, token.Resource).Unwrap()

	id := parseId(lexer).Unwrap()

	var methods []ast.ResourceMethod

	tok := peek(lexer).Unwrap()
	switch tok.Type {
	case token.Semicolon:
		expect(lexer, token.Semicolon).Unwrap()
	case token.OpenBrace:
		expect(lexer, token.OpenBrace).Unwrap()

		tok := peek(lexer).Unwrap()
		for tok.Type != token.CloseBrace {
			method := parseResourceMethod(lexer).Unwrap()
			methods = append(methods, method)
			tok = peek(lexer).Unwrap()
		}
		expect(lexer, token.CloseBrace).Unwrap()
	}

	return result.Ok(ast.Resource{
		ID:      id,
		Methods: methods,
	})
}

func parseResourceMethod(lexer *lex.Lexer) (res types.Result[ast.ResourceMethod]) {

	var resourceMethod ast.ResourceMethod

	// resource-method ::= 'constructor' param-list ';'
	if eat(lexer, token.Constructor).Unwrap() {
		parameters := parseParameters(lexer).Unwrap()
		expect(lexer, token.Semicolon).Unwrap()
		resourceMethod = &ast.Constructor{
			ParameterList: parameters,
		}
		return result.Ok(resourceMethod)
	}

	// the resource-method with func-item overlaps with the resource item static for the first two tokens
	// create a clone of the lexer and commit the changes if the keyword 'static' occurs after the colon
	// otherwise throw the clone away and parse as a func item
	clone := lexer.Clone()

	// id
	id := parseId(clone).Unwrap()

	// ';'
	expect(clone, token.Colon).Unwrap()

	if !eat(clone, token.Static).Unwrap() {
		// resource-method ::= func-item
		// func-item ::= id ':' func-type ';'
		// throw away the clone
		funcItem := parseFuncItem(lexer).Unwrap()
		resourceMethod = ast.Method{
			Func: funcItem,
		}
	} else {
		// resource-method ::= id ':' 'static' func-type ';'
		*lexer = *clone
		funcType := parseFunc(lexer).Unwrap()
		expect(lexer, token.Semicolon).Unwrap()
		resourceMethod = ast.Static{
			ID:       id,
			FuncType: funcType,
		}
	}
	return result.Ok(resourceMethod)
}

func parseUse(lexer *lex.Lexer) (res types.Result[*ast.Use]) {
	defer handle.Error(&res)

	// 'use'
	expect(lexer, token.Use).Unwrap()

	// use-path
	from := parseUsePath(lexer).Unwrap()

	// .
	expect(lexer, token.Period).Unwrap()
	names := parseItemList[ast.UseName](
		lexer,
		token.OpenBrace,
		token.CloseBrace,
		parseUseName).Unwrap()

	// ;
	expect(lexer, token.Semicolon).Unwrap()

	return result.Ok(&ast.Use{
		From:  from,
		Names: names,
	})
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
	usePath := &ast.UsePath{
		Id: "string",
		Package: struct {
			Id   *ast.PackageDeclaration
			Name string
		}{
			Id: &ast.PackageDeclaration{
				Namespace: namespace,
				Name:      pkgName,
				Version:   version,
			},
			Name: name,
		},
	}
	return result.Ok(usePath)
}

func parseOptionalVersion(lexer *lex.Lexer) (res types.Result[types.Option[ast.Version]]) {
	return result.Errorf[types.Option[ast.Version]]("not implemented")
}

func parseFuncItem(lexer *lex.Lexer) (res types.Result[*ast.FuncItem]) {
	defer handle.Error(&res)

	// func-item ::= id ':' func-type ';'
	id := parseId(lexer).Unwrap()
	expect(lexer, token.Colon).Unwrap()
	funcType := parseFunc(lexer).Unwrap()
	expect(lexer, token.Semicolon).Unwrap()

	return result.Ok(&ast.FuncItem{
		ID:       id,
		FuncType: funcType,
	})
}

func parseFunc(lexer *lex.Lexer) (res types.Result[*ast.FuncType]) {

	defer handle.Error(&res)

	expect(lexer, token.Func).Unwrap()

	parameters := parseParameters(lexer).Unwrap()
	results := &ast.ResultList{}
	if eat(lexer, token.RightArrow).Unwrap() {
		tok := peek(lexer).Unwrap()
		if tok.Type == token.OpenParen {
			results.Named = parseParameters(lexer).Unwrap()
		} else {
			results.Anonymous = parseType(lexer).Unwrap()
		}
	} else {
		results.Named = nil // ? []ast.Parameter{}
	}

	return result.Ok(&ast.FuncType{
		Params:  parameters,
		Results: results,
	})
}

func parseParameters(lexer *lex.Lexer) (res types.Result[[]ast.Parameter]) {
	var parameters []ast.Parameter

	expect(lexer, token.OpenParen).Unwrap()

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

func parseTypeItem(lexer *lex.Lexer) (res types.Result[*ast.TypeItem]) {
	defer handle.Error(&res)

	expect(lexer, token.Type).Unwrap()

	id := parseId(lexer).Unwrap()
	expect(lexer, token.Equal).Unwrap()

	ty := parseType(lexer).Unwrap()
	expect(lexer, token.Semicolon).Unwrap()

	return result.Ok(&ast.TypeItem{
		ID:   id,
		Type: ty,
	})
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

	future := &ast.Future{
		ItemType: option.None[ast.Type](),
	}

	if eat(lexer, token.Less).Unwrap() {
		ty := parseType(lexer).Unwrap()
		future.ItemType = option.Some(ty)
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

	itemType := peek(lexer).Unwrap()
	var worldItem ast.WorldItem
	switch itemType.Type {
	case token.Export:
		worldItem = parseExport(lexer).Unwrap()
	case token.Import:
		worldItem = parseImport(lexer).Unwrap()
	case token.Use:
		worldItem = parseUse(lexer).Unwrap()
	case token.Type:
		worldItem = parseTypeDef(lexer).Unwrap()
	case token.Record:
		worldItem = parseRecord(lexer).Unwrap()
	case token.Variant:
		worldItem = parseVariant(lexer).Unwrap()
	case token.Resource:
		worldItem = parseResource(lexer).Unwrap()
	case token.Include:
		worldItem = parseInclude(lexer).Unwrap()

	default:
		return result.Errorf[ast.WorldItem]("%w : Unrecognized world item '%s'. Expected (export, import, resource, use, type, include). Found token.%v",
			parseErrorFromLexer(lexer),
			itemType.Capture,
			itemType.Type)
	}
	return result.Ok(worldItem)
}

func parseExport(lexer *lex.Lexer) (res types.Result[*ast.Export]) {
	defer handle.Error(&res)
	expect(lexer, token.Export).Unwrap()
	ty := parseExternType(lexer).Unwrap()
	return result.Ok(&ast.Export{
		ExternType: ty,
	})
}

func parseImport(lexer *lex.Lexer) (res types.Result[*ast.Import]) {
	defer handle.Error(&res)
	expect(lexer, token.Import).Unwrap()
	ty := parseExternType(lexer).Unwrap()
	return result.Ok(&ast.Import{
		ExternType: ty,
	})
}

func parseInclude(lexer *lex.Lexer) (res types.Result[ast.WorldItem]) {
	defer handle.Error(&res)

	expect(lexer, token.Include).Unwrap()

	// include-item = 'include' use-path ';'
	// include-item = 'include' use-path 'with' '{' include-names-list '}'
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
	} else {
		expect(lexer, token.Semicolon).Unwrap()
	}

	return result.Ok[ast.WorldItem](include)
}

func parseExternType(lexer *lex.Lexer) (res types.Result[ast.ExternType]) {
	defer handle.Error(&res)

	// There is some ambiguity in how functions, interface and use path import/export overlap.
	// Clone the lexer and try to make progress with interface and function import/export
	// if successful, apply the clone's progress to the input lexer and continue parsing
	// if failed, try to parse using a use path
	clone := lexer.Clone()
	id := parseId(clone).Unwrap()
	if !eat(clone, token.Colon).Unwrap() {

		usePath := parseUsePath(lexer).Unwrap()
		expect(lexer, token.Semicolon).Unwrap()

		return result.Ok[ast.ExternType](&ast.ExternTypeUsePath{
			UsePath: usePath,
		})
	}

	tok := peek(clone).Unwrap()
	switch tok.Type {
	case token.Func:
		*lexer = *clone
		function := parseFunc(lexer).Unwrap()

		expect(lexer, token.Semicolon).Unwrap()
		return result.Ok[ast.ExternType](&ast.ExternTypeFunc{
			ID:   id,
			Func: function,
		})
	case token.Interface:
		*lexer = *clone

		expect(lexer, token.Interface).Unwrap()
		return result.Ok[ast.ExternType](&ast.ExternTypeInterface{
			ID:             id,
			InterfaceItems: parseInterfaceItems(lexer).Unwrap(),
		})
	}

	return result.Errorf[ast.ExternType]("unable to parse ExternType %w", parseError(tok))
}

// record-item ::= 'record' id '{' record-fields '}'
// record-fields ::= record-field | record-field ',' record-fields?
// record-field ::= id ':' ty
func parseRecord(lexer *lex.Lexer) (res types.Result[*ast.Record]) {
	defer handle.Error(&res)

	// 'record'
	expect(lexer, token.Record).Unwrap()

	name := parseId(lexer).Unwrap()

	fields := parseItemList(
		lexer,
		token.OpenBrace, token.CloseBrace,
		parseRecordField).Unwrap()

	return result.Ok(&ast.Record{
		ID:     name,
		Fields: fields,
	})
}

func parseRecordField(lexer *lex.Lexer) (res types.Result[ast.Field]) {
	defer handle.Error(&res)
	id := parseId(lexer).Unwrap()
	expect(lexer, token.Colon).Unwrap()
	ty := parseType(lexer).Unwrap()
	return result.Ok(ast.Field{Name: id, Type: ty})
}

// flags-items ::= 'flags' id '{' flags-fields '}'
// flags-fields ::= id  | id ',' flags-fields?
func parseFlags(lexer *lex.Lexer) (res types.Result[*ast.Flags]) {
	defer handle.Error(&res)

	// 'flags'
	expect(lexer, token.Flags).Unwrap()

	name := parseId(lexer).Unwrap()
	flagList := parseItemList(lexer, token.OpenBrace, token.CloseBrace, func(l *lex.Lexer) types.Result[ast.Flag] {
		id := parseId(lexer).Unwrap()
		return result.Ok(ast.Flag{
			Id: id,
		})
	}).Unwrap()

	return result.Ok(&ast.Flags{
		ID:    name,
		Flags: flagList,
	})
}

// variant-items ::= 'variant' id '{' variant-cases '}'
// variant-cases ::= variant-case | variant-case ',' variant-cases?
// variant-case ::= id | id '(' ty ')'
func parseVariant(lexer *lex.Lexer) (res types.Result[*ast.Variant]) {
	defer handle.Error(&res)

	// 'variant'
	expect(lexer, token.Variant).Unwrap()

	// id
	name := parseId(lexer).Unwrap()

	// '{' variant-cases '}'
	cases := parseItemList(
		lexer,
		token.OpenBrace, token.CloseBrace,
		parseVariantCase).Unwrap()

	return result.Ok(
		&ast.Variant{
			ID:    name,
			Cases: cases,
		})
}

func parseVariantCase(lexer *lex.Lexer) (res types.Result[ast.Case]) {
	defer handle.Error(&res)
	name := parseId(lexer).Unwrap()
	c := &ast.Case{
		Name: name,
		Type: option.None[ast.Type](),
	}
	if eat(lexer, token.OpenParen).Unwrap() {
		ty := parseType(lexer).Unwrap()
		expect(lexer, token.CloseParen).Unwrap()
		c.Type = option.Some(ty)
	}
	return result.Ok(*c)
}

// enum-items ::= 'enum' id '{' enum-cases '}'
// enum-cases ::= id | id ',' enum-cases?
func parseEnum(lexer *lex.Lexer) (res types.Result[*ast.Enum]) {
	defer handle.Error(&res)

	// 'enum'
	expect(lexer, token.Enum).Unwrap()

	// id
	id := parseId(lexer).Unwrap()

	// '{' enum-cases '}'
	cases := parseItemList(
		lexer,
		token.OpenBrace, token.CloseBrace,
		parseEnumCase).Unwrap()

	return result.Ok(&ast.Enum{
		Cases: cases,
		ID:    id,
	})
}

func parseEnumCase(lexer *lex.Lexer) (res types.Result[ast.EnumCase]) {
	defer handle.Error(&res)
	id := parseId(lexer).Unwrap()
	return result.Ok(ast.EnumCase{
		Name: id,
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

func parseInteger(lexer *lex.Lexer) (res types.Result[int64]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	switch tok.Type {
	case token.Integer:
		i, err := strconv.ParseInt(tok.Capture, 0, 32)
		return result.New(i, err)
	default:
		return result.Errorf[int64]("%w: found value '%s', type '%v' but expected (integer) ", parseError(tok), tok.Capture, tok.Type)
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
		tok := p.Unwrap()

		if tok.Type != token.Whitespace &&
			tok.Type != token.BlockComment &&
			tok.Type != token.LineComment {
			return p
		}

		// consume ignore tokens
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
