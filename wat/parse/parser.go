package parse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patrickhuber/go-types"
	"github.com/patrickhuber/go-types/handle"
	"github.com/patrickhuber/go-types/option"
	"github.com/patrickhuber/go-types/result"
	"github.com/patrickhuber/go-wasm/wat/ast"
	"github.com/patrickhuber/go-wasm/wat/lex"

	"github.com/patrickhuber/go-wasm/wat/token"
)

func Parse(lexer *lex.Lexer) (ast.Ast, error) {
	return parse(lexer).Deconstruct()
}

func Peek(lexer *lex.Lexer) (*token.Token, error) {
	return peek(lexer).Deconstruct()
}

func parse(lexer *lex.Lexer) (res types.Result[ast.Ast]) {
	defer handle.Error(&res)
	expect(lexer, token.OpenParen).Unwrap()
	tok := peek(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Errorf[ast.Ast]("%w : unrecognized token", parseError(tok))
	}
	var root ast.Ast
	switch tok.Capture {
	case "module":
		root = parseModule(lexer).Unwrap()
	case "component":
		root = parseComponent(lexer).Unwrap()
	default:
		return result.Errorf[ast.Ast](
			"%w : expected module, component but found '%s'",
			parseError(tok), tok.Capture)
	}
	expect(lexer, token.CloseParen).Unwrap()
	return result.Ok(root)
}

func parseModule(lexer *lex.Lexer) (res types.Result[*ast.Module]) {
	defer handle.Error(&res)
	expect(lexer, token.Reserved).Unwrap()

	m := &ast.Module{}
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := peek(lexer).Unwrap()
		switch tok.Capture {
		case "func":
			f := parseFunc(lexer).Unwrap()
			m.Functions = append(m.Functions, *f)
		case "type":
			t := parseType(lexer).Unwrap()
			m.Types = append(m.Types, t)
		case "table":
			t := parseTable(lexer).Unwrap()
			m.Tables = append(m.Tables, t)
		case "global":
			g := parseGlobal(lexer).Unwrap()
			m.Globals = append(m.Globals, g)
		case "memory":
			mem := parseMemory(lexer).Unwrap()
			m.Memory = append(m.Memory, mem)
		default:
			return result.Errorf[*ast.Module]("unrecognized module section '%s'", tok.Capture)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok(m)
}

func parseComponent(lexer *lex.Lexer) (res types.Result[*ast.Component]) {
	defer handle.Error(&res)
	expect(lexer, token.Reserved).Unwrap()
	return result.Ok(&ast.Component{})
}

func parseFunc(lexer *lex.Lexer) (res types.Result[*ast.Function]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "func").Unwrap()

	function := &ast.Function{}
	if tok := peek(lexer).Unwrap(); tok.Type == token.Id {
		id := next(lexer).Unwrap()
		if id.Type == token.Id {
			function.ID = option.Some(tok.Capture)
		} else {
			function.ID = option.None[string]()
		}
	}
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := peek(lexer).Unwrap()
		switch tok.Capture {
		case "local":
			local := parseLocal(lexer).Unwrap()
			function.Locals = append(function.Locals, local)
		case "param":
			param := parseParameter(lexer).Unwrap()
			function.Parameters = append(function.Parameters, *param)
		case "result":
			result := parseResult(lexer).Unwrap()
			function.Results = append(function.Results, *result)
		case "export":
			export := parseExport(lexer).Unwrap()
			function.Exports = append(function.Exports, export)
		case "import":
			_ = parseImport(lexer).Unwrap()
		default:
			inst := parseInstruction(lexer).Unwrap()
			function.Instructions = append(function.Instructions, inst)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}
	for {
		tok := peek(lexer).Unwrap()
		if tok.Type == token.CloseParen {
			break
		}
		instruction := parseInstruction(lexer).Unwrap()
		function.Instructions = append(function.Instructions, instruction)
	}
	return result.Ok(function)
}

func parseLocal(lexer *lex.Lexer) (res types.Result[ast.Local]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "local").Unwrap()
	ty := parseValType(lexer).Unwrap()
	return result.Ok(ast.Local{
		Type: ty,
	})
}

func parseParameter(lexer *lex.Lexer) (res types.Result[*ast.Parameter]) {
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "param").Unwrap()

	id := parseOptionalId(lexer).Unwrap()
	var types []ast.ValType
	if id.IsSome() {
		types = append(types, parseValType(lexer).Unwrap())
	} else {
		tok := peek(lexer).Unwrap()
		for tok.Type == token.Reserved {
			ty := parseValType(lexer).Unwrap()
			types = append(types, ty)
			tok = peek(lexer).Unwrap()
		}
	}

	return result.Ok(&ast.Parameter{
		ID:    id,
		Types: types,
	})
}

func parseResult(lexer *lex.Lexer) (res types.Result[*ast.Result]) {
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "result").Unwrap()

	ty := parseValType(lexer).Unwrap()
	return result.Ok(&ast.Result{
		Types: []ast.ValType{ty},
	})
}

func parseType(lexer *lex.Lexer) (res types.Result[ast.Type]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "type").Unwrap()

	id := parseOptionalId(lexer).Unwrap()

	var funcType ast.FuncType
	if eat(lexer, token.OpenParen).Unwrap() {
		funcType = parseFuncType(lexer).Unwrap()
		expect(lexer, token.CloseParen).Unwrap()
	}

	return result.Ok(ast.Type{
		ID:       id,
		FuncType: funcType,
	})
}

func parseFuncType(lexer *lex.Lexer) (res types.Result[ast.FuncType]) {
	defer handle.Error(&res)

	expectValue(lexer, token.Reserved, "func").Unwrap()
	var parameters []ast.Parameter
	var results []ast.Result
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := peek(lexer).Unwrap()
		switch tok.Capture {
		case "param":
			parameter := parseParameter(lexer).Unwrap()
			parameters = append(parameters, *parameter)
		case "result":
			result := parseResult(lexer).Unwrap()
			results = append(results, *result)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}

	return result.Ok(ast.FuncType{
		Parameters: parameters,
		Results:    results,
	})
}

func parseTable(lexer *lex.Lexer) (res types.Result[ast.Table]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "table").Unwrap()
	id := parseOptionalId(lexer).Unwrap()
	tableType := parseTableType(lexer).Unwrap()
	var elements []ast.Element
	for eat(lexer, token.OpenParen).Unwrap() {
		tok := next(lexer).Unwrap()
		if tok.Type != token.Reserved {
			return result.Errorf[ast.Table]("%w", parseError(tok))
		}
		switch tok.Capture {
		case "elem":
			id := parseId(lexer).Unwrap()
			elements = append(elements, ast.Element{
				ID: id,
			})
		default:
			return result.Errorf[ast.Table]("%w", parseError(tok))
		}
		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok(ast.Table{
		ID:        id,
		TableType: tableType,
		Elements:  elements,
	})
}

func parseTableType(lexer *lex.Lexer) (res types.Result[ast.TableType]) {
	defer handle.Error(&res)

	tok := peek(lexer).Unwrap()
	var limits ast.Limits
	if tok.Type == token.Integer {
		limits = parseLimits(lexer).Unwrap()
	}
	tok = next(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Errorf[ast.TableType]("expected 'externref' or 'funcref': %w", parseError(tok))
	}
	var refType ast.RefType
	switch tok.Capture {
	case "externref":
		refType = ast.ExternRef{}
	case "funcref":
		refType = ast.FuncRef{}
	}

	return result.Ok(ast.TableType{
		Limits:  limits,
		RefType: refType,
	})
}

func parseGlobal(lexer *lex.Lexer) (res types.Result[ast.Global]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "global").Unwrap()
	id := parseOptionalId(lexer).Unwrap()
	globalType := parseGlobalType(lexer).Unwrap()
	instructions := parseInstructions(lexer).Unwrap()
	return result.Ok(ast.Global{
		ID:           id,
		Type:         globalType,
		Instructions: instructions,
	})
}

func parseGlobalType(lexer *lex.Lexer) (res types.Result[ast.GlobalType]) {
	defer handle.Error(&res)
	var mutable bool
	var valType ast.ValType
	if eat(lexer, token.OpenParen).Unwrap() {
		expectValue(lexer, token.Reserved, "mut").Unwrap()
		mutable = true
		valType = parseValType(lexer).Unwrap()
		expect(lexer, token.CloseParen).Unwrap()
	} else {
		valType = parseValType(lexer).Unwrap()
	}
	return result.Ok(ast.GlobalType{
		Type:    valType,
		Mutable: mutable,
	})
}

func parseMemory(lexer *lex.Lexer) (res types.Result[ast.Memory]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "memory").Unwrap()
	return result.Ok(ast.Memory{
		ID:     parseOptionalId(lexer).Unwrap(),
		Limits: parseLimits(lexer).Unwrap(),
	})
}

func parseLimits(lexer *lex.Lexer) (res types.Result[ast.Limits]) {
	defer handle.Error(&res)
	min := parseInt32(lexer).Unwrap()
	max := option.None[uint32]()
	tok := peek(lexer).Unwrap()
	if tok.Type == token.Integer {
		n := parseInt32(lexer).Unwrap()
		max = option.Some(uint32(n))
	}
	return result.Ok(ast.Limits{
		Min: uint32(min),
		Max: max,
	})
}

func parseExport(lexer *lex.Lexer) (res types.Result[ast.InlineExport]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "export").Unwrap()
	name := parseString(lexer).Unwrap()
	return result.Ok(ast.InlineExport{
		Name: name,
	})
}

func parseImport(lexer *lex.Lexer) (res types.Result[ast.InlineImport]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "import").Unwrap()
	name := parseString(lexer).Unwrap()
	alias := parseString(lexer).Unwrap()
	return result.Ok(ast.InlineImport{
		Module: name,
		Field:  alias,
	})
}

func parseInstructions(lexer *lex.Lexer) (res types.Result[[]ast.Instruction]) {
	defer handle.Error(&res)

	var instructions []ast.Instruction
	for eat(lexer, token.OpenParen).Unwrap() {
		instruction := parseInstruction(lexer).Unwrap()
		instructions = append(instructions, instruction)
		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok(instructions)
}

func parseInstruction(lexer *lex.Lexer) (res types.Result[ast.Instruction]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != token.Reserved {
		return result.Error[ast.Instruction](parseError(tok))
	}
	var inst ast.Instruction
	switch tok.Capture {
	case "br":
		inst = ast.Br{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "br_if":
		inst = ast.BrIf{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "br_table":
		inst = ast.BrTable{
			Indicies: []ast.Index{
				parseIndex(lexer).Unwrap(),
			},
		}
	case "return":
		inst = ast.Return{}
	case "call":
		inst = ast.Call{
			Index: parseIndex(lexer).Unwrap(),
		}

	case "drop":
		inst = ast.Drop{}
	case "select":
		inst = ast.Select{}
	case "local.get":
		inst = ast.LocalGet{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "local.set":
		inst = ast.LocalSet{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "local.tee":
		inst = ast.LocalTee{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "global.get":
		inst = ast.GlobalGet{
			Index: parseIndex(lexer).Unwrap(),
		}
	case "global.set":
		inst = ast.GlobalSet{
			Index: parseIndex(lexer).Unwrap(),
		}

	// memory instructions
	case "memory.grow":
		inst = ast.MemoryGrow{}
	case "i32.load":
		inst = ast.I32Load{}
	case "i32.store":
		inst = ast.I32Store{}

	// numeric instructions
	case "i32.const":
		inst = ast.I32Const{
			Value: parseInt32(lexer).Unwrap(),
		}
	case "i64.const":
		inst = ast.I64Const{
			Value: parseInt64(lexer).Unwrap(),
		}
	case "f32.const":
		inst = ast.F32Const{
			Value: parseFloat32(lexer).Unwrap(),
		}
	case "f64.const":
		inst = ast.F64Const{
			Value: parseFloat64(lexer).Unwrap(),
		}
	case "i32.add":
		inst = ast.I32Add{}
	case "i32.sub":
		inst = ast.I32Sub{}
	case "i32.mul":
		inst = ast.I32Mul{}
	case "i32.div_s":
		inst = ast.I32DivS{}
	case "i32.div_u":
		inst = ast.I32DivU{}
	case "i32.rem_s":
	case "i32.rem_u":
	case "i32.and":
	case "i32.or":
	case "i32.xor":
	case "i32.shl":
	case "i32.shr_s":
	case "i32.shr_u":
	case "i32.rotl":
	case "i32.rotr":
	case "i32.clz":
	case "i32.ctz":
	case "i32.popcnt":
	case "i32.extend8_s":
	case "i32.extend16_s":
	case "i32.eqz":
		inst = ast.I32Eqz{}
	case "i32.eq":
	case "i32.ne":
	case "i32.lt_s":
	case "i32.lt_u":
	case "i32.le_s":
	case "i32.le_u":
	case "i32.gt_s":
	case "i32.gt_u":
	case "i32.ge_s":
	case "i32.ge_u":
	case "block":
		inst = parseBlock(lexer).Unwrap()
	case "loop":
		inst = parseLoop(lexer).Unwrap()
	case "if":
		inst = parseIf(lexer).Unwrap()

	case "call_indirect":
		inst = parseCallIndirect(lexer).Unwrap()
	default:
		return result.Errorf[ast.Instruction]("%w : error parsing instruction. Unrecognized instruction %v : %s", parseError(tok), tok.Type, tok.Capture)
	}
	peekTok := peek(lexer).Unwrap()
	if peekTok.Type != token.OpenParen {
		return result.Ok(inst)
	}
	folded := ast.Folded{
		Instruction: inst,
	}
	for eat(lexer, token.OpenParen).Unwrap() {
		inst = parseInstruction(lexer).Unwrap()
		folded.Parameters = append(folded.Parameters, inst)
		expect(lexer, token.CloseParen).Unwrap()
	}
	return result.Ok[ast.Instruction](folded)
}

func parseIndex(lexer *lex.Lexer) (res types.Result[ast.Index]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	var index ast.Index
	switch tok.Type {
	case token.Integer:
		i := result.New(strconv.Atoi(tok.Capture)).Unwrap()
		index = &ast.RawIndex{
			Index: uint32(i),
		}
	case token.Id:
		index = &ast.IDIndex{
			ID: tok.Capture,
		}
	default:
		return result.Errorf[ast.Index]("%w : error parsing index", parseError(tok))
	}
	return result.Ok(index)
}

func parseValType(lexer *lex.Lexer) (res types.Result[ast.ValType]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()

	var ty ast.ValType
	// todo: enhance the lexer to parse these as tokens
	switch tok.Capture {
	case "i32":
		ty = ast.I32{}
	case "i64":
		ty = ast.I64{}
	case "f32":
		ty = ast.F32{}
	case "f64":
		ty = ast.F64{}
	default:
		return result.Errorf[ast.ValType]("%w : error parsing type. expected (i32, i64, f32, f64) but found %s", parseError(tok), tok.Capture)
	}
	return result.Ok(ty)
}

func parseBlock(lexer *lex.Lexer) (res types.Result[ast.Block]) {
	defer handle.Error(&res)
	tok := peek(lexer).Unwrap()

	name := option.None[string]()
	if tok.Type == token.Id {
		n := next(lexer).Unwrap().Capture
		name = option.Some(n)
	}

	blockType := parseBlockType(lexer).Unwrap()
	instructions := parseInstructions(lexer).Unwrap()

	return result.Ok(ast.Block{
		Name:         name,
		BlockType:    blockType,
		Instructions: instructions,
	})
}

func parseLoop(lexer *lex.Lexer) (res types.Result[ast.Loop]) {
	defer handle.Error(&res)
	tok := peek(lexer).Unwrap()

	name := option.None[string]()
	if tok.Type == token.Id {
		n := next(lexer).Unwrap().Capture
		name = option.Some(n)
	}

	blockType := parseBlockType(lexer).Unwrap()
	instructions := parseInstructions(lexer).Unwrap()

	return result.Ok(ast.Loop{
		Name:         name,
		BlockType:    blockType,
		Instructions: instructions,
	})
}

func parseIf(lexer *lex.Lexer) (res types.Result[ast.If]) {
	defer handle.Error(&res)
	tok := peek(lexer).Unwrap()

	name := option.None[string]()
	if tok.Type == token.Id {
		n := next(lexer).Unwrap().Capture
		name = option.Some(n)
	}

	var results []ast.Result
	_else := option.None[ast.Else]()
	var then ast.Then
	var instructions []ast.Instruction
	for eat(lexer, token.OpenParen).Unwrap() {
		tok = peek(lexer).Unwrap()
		switch tok.Capture {
		case "result":
			res := parseResult(lexer).Unwrap()
			results = append(results, *res)
		case "then":
			then = parseThen(lexer).Unwrap()
		case "else":
			e := parseElse(lexer).Unwrap()
			_else = option.Some(e)
		default:
			i := parseInstruction(lexer).Unwrap()
			instructions = append(instructions, i)
		}
		expect(lexer, token.CloseParen).Unwrap()
	}

	return result.Ok(ast.If{
		Name: name,
		BlockType: ast.BlockType{
			Results: results,
		},
		Clause: instructions,
		Then:   then,
		Else:   _else,
	})
}

func parseElse(lexer *lex.Lexer) (res types.Result[ast.Else]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "else").Unwrap()
	return result.Ok(ast.Else{
		Instructions: parseInstructions(lexer).Unwrap(),
	})
}

func parseThen(lexer *lex.Lexer) (res types.Result[ast.Then]) {
	defer handle.Error(&res)
	expectValue(lexer, token.Reserved, "then").Unwrap()
	return result.Ok(ast.Then{
		Instructions: parseInstructions(lexer).Unwrap(),
	})
}

func parseCallIndirect(lexer *lex.Lexer) (res types.Result[ast.CallIndirect]) {
	defer handle.Error(&res)

	typeUse := parseTypeUse(lexer).Unwrap()
	return result.Ok(ast.CallIndirect{
		Type: typeUse,
	})
}

func parseTypeUse(lexer *lex.Lexer) (res types.Result[ast.TypeUse]) {
	defer handle.Error(&res)
	expect(lexer, token.OpenParen).Unwrap()
	expectValue(lexer, token.Reserved, "type").Unwrap()
	index := parseId(lexer).Unwrap()
	expect(lexer, token.CloseParen).Unwrap()
	return result.Ok(ast.TypeUse{
		Index: index,
	})
}

func parseBlockType(lexer *lex.Lexer) (res types.Result[ast.BlockType]) {
	defer handle.Error(&res)

	// create a clone of the lexer to check for any results
	clone := lexer.Clone()

	// block_type = ( result <val_type>* )*
	var results []ast.Result
	for eat(clone, token.OpenParen).Unwrap() {
		tok := peek(clone).Unwrap()

		// if we are not at a result, this is an instruction so roll back
		if tok.Capture != "result" {
			break
		}

		// merge the lexer back because we know we are parsing a result
		*lexer = *clone

		expectValue(lexer, token.Reserved, "result").Unwrap()

		var types []ast.ValType
		for tok := peek(lexer).Unwrap(); tok.Type != token.CloseParen; tok = peek(lexer).Unwrap() {
			ty := parseValType(lexer).Unwrap()
			types = append(types, ty)
		}
		result := ast.Result{
			Types: types,
		}
		results = append(results, result)
		expect(lexer, token.CloseParen).Unwrap()

		// create a new clone and continue parsing
		clone = lexer.Clone()
	}
	return result.Ok(ast.BlockType{
		Results: results,
	})
}

func parseOptionalId(lexer *lex.Lexer) (res types.Result[types.Option[string]]) {
	tok := peek(lexer).Unwrap()
	if tok.Type == token.Id {
		id := parseId(lexer).Unwrap()
		return result.Ok(option.Some(id))
	}
	return result.Ok(option.None[string]())
}

func parseId(lexer *lex.Lexer) (res types.Result[string]) {
	tok := next(lexer).Unwrap()
	if tok.Type != token.Id {
		return result.Errorf[string]("%w", parseError(tok))
	}
	return result.Ok(tok.Capture)
}

func ParseString(lexer *lex.Lexer) (string, error) {
	return parseString(lexer).Deconstruct()
}

func parseString(lexer *lex.Lexer) (res types.Result[string]) {
	tok := next(lexer).Unwrap()
	if tok.Type != token.String {
		return result.Errorf[string]("%w", parseError(tok))
	}
	return result.Ok(strings.Trim(tok.Capture, "\""))
}

func ParseInt32(lexer *lex.Lexer) (int32, error) {
	return parseInt32(lexer).Deconstruct()
}

func parseInt32(lexer *lex.Lexer) (res types.Result[int32]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	if tok.Type != token.Integer {
		return result.Errorf[int32]("expected integer %w", parseError(tok))
	}
	if strings.HasPrefix(tok.Capture, "-") {
		i, err := strconv.ParseInt(tok.Capture, 0, 32)
		return result.New(int32(i), err)
	}
	i, err := strconv.ParseUint(tok.Capture, 0, 32)
	return result.New(int32(i), err)
}

func ParseInt64(lexer *lex.Lexer) (int64, error) {
	return parseInt64(lexer).Deconstruct()
}

func parseInt64(lexer *lex.Lexer) (res types.Result[int64]) {
	defer handle.Error(&res)

	tok := next(lexer).Unwrap()
	if tok.Type != token.Integer {
		return result.Errorf[int64]("expected integer %w", parseError(tok))
	}
	if strings.HasPrefix(tok.Capture, "-") {
		i, err := strconv.ParseInt(tok.Capture, 0, 32)
		return result.New(int64(i), err)
	}
	i, err := strconv.ParseUint(tok.Capture, 0, 32)
	return result.New(int64(i), err)
}

func ParseFloat32(lexer *lex.Lexer) (float32, error) {
	return parseFloat32(lexer).Deconstruct()
}

func parseFloat32(lexer *lex.Lexer) (res types.Result[float32]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != token.Float {
		return result.Errorf[float32]("expected float %w", parseError(tok))
	}
	f, err := strconv.ParseFloat(tok.Capture, 32)
	return result.New(float32(f), err)
}

func ParseFloat64(lexer *lex.Lexer) (float64, error) {
	return parseFloat64(lexer).Deconstruct()
}

func parseFloat64(lexer *lex.Lexer) (res types.Result[float64]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != token.Float {
		return result.Errorf[float64]("expected float %w", parseError(tok))
	}
	f, err := strconv.ParseFloat(tok.Capture, 32)
	return result.New(f, err)
}

func eat(lexer *lex.Lexer, ty token.Type) (res types.Result[bool]) {
	defer handle.Error(&res)

	tok := peek(lexer).Unwrap()
	if tok.Type != ty {
		return result.Ok(false)
	}

	expect(lexer, ty).Unwrap()
	return result.Ok(true)
}

func expectValue(lexer *lex.Lexer, ty token.Type, capture string) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type != ty {
		return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), ty, tok.Type)
	}
	if tok.Capture != capture {
		return result.Errorf[any]("%w. expected type:'%v' value:'%s' but found type:'%v' value:'%s'",
			parseError(tok), ty, capture, tok.Type, tok.Capture)
	}
	return result.Ok[any](nil)
}

func expect(lexer *lex.Lexer, ty token.Type) (res types.Result[any]) {
	defer handle.Error(&res)
	tok := next(lexer).Unwrap()
	if tok.Type == ty {
		return result.Ok[any](nil)
	}
	return result.Errorf[any]("%w. expected '%v' but found '%v'", parseError(tok), ty, tok.Type)
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

func peek(lexer *lex.Lexer) (res types.Result[*token.Token]) {
	defer handle.Error(&res)
	for {
		p := result.New(lexer.Peek())
		r := p.Unwrap()
		switch r.Type {
		case token.Whitespace:
		case token.BlockComment:
		case token.LineComment:
		default:
			return p
		}
		// consume whitespace, block comment and line comment
		_ = result.New(lexer.Next()).Unwrap()
	}
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
