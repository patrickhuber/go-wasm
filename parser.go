package wasm

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(lexer Lexer) (*Module, error) {
	return NewParser(lexer).Parse()
}

func ParseString(input string) (*Module, error) {
	lexer := NewLexer(input)
	return Parse(lexer)
}

type Parser interface {
	Parse() (*Module, error)
}

type parser struct {
	lexer Lexer
}

func NewParser(lexer Lexer) Parser {

	return &parser{
		lexer: lexer,
	}
}

func (p *parser) Parse() (*Module, error) {
	return p.ParseModule()
}

func (p *parser) ParseModule() (*Module, error) {

	err := p.ExpectToken(OpenParen)
	if err != nil {
		return nil, err
	}

	err = p.ExpectString("module")
	if err != nil {
		return nil, err
	}

	module := &Module{}
	for {
		section, ok, err := p.TryParseSection()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		if section.Function != nil {
			module.Functions = append(module.Functions, *section.Function)
		}
		if section.Memory != nil {
			module.Memory = append(module.Memory, *section.Memory)
		}
	}
	err = p.ExpectToken(CloseParen)
	if err != nil {
		return nil, err
	}

	return module, err
}

func (p *parser) TryParseSection() (*Section, bool, error) {
	peek, err := p.peekToken()
	if err != nil {
		return nil, false, err
	}

	if peek.Type == CloseParen {
		return nil, false, nil
	}

	section, err := p.ParseSection()
	if err != nil {
		return nil, false, err
	}

	return section, true, err

}

func (p *parser) ParseSection() (*Section, error) {
	err := p.ExpectToken(OpenParen)
	if err != nil {
		return nil, err
	}

	tok, err := p.ParseString()
	if err != nil {
		return nil, err
	}

	var section *Section

	switch tok.Capture {
	case "func":
		function, err := p.ParseFunction()
		if err != nil {
			return nil, err
		}
		section = &Section{
			Function: function,
		}
	case "memory":
		memory, err := p.ParseMemory()
		if err != nil {
			return nil, err
		}
		section = &Section{
			Memory: memory,
		}
	default:
		return nil, p.parseError(tok, fmt.Errorf("unexpected token %s found", tok.Type))
	}
	return section, p.ExpectToken(CloseParen)
}

func (p *parser) ParseFunction() (*Function, error) {
	function, err := p.ParseSignature()
	if err != nil {
		return nil, err
	}
	instructions, err := p.ParseInstructions()
	if err != nil {
		return nil, err
	}
	function.Instructions = instructions
	return function, nil
}

func (p *parser) ParseMemory() (*Memory, error) {
	memory := &Memory{
		Limits: Limits{},
	}
	identifier, ok, err := p.TryParseIdentifier()
	if err != nil {
		return nil, err
	}
	if ok {
		memory.ID = identifier
	}

	min, err := p.ParseUInt32()
	if err != nil {
		return nil, err
	}
	memory.Limits.Min = min

	peek, err := p.peekToken()
	if err != nil {
		return nil, err
	}

	if peek.Type == OpenParen || peek.Type == CloseParen {
		return memory, nil
	}

	max, err := p.ParseUInt32()
	if err != nil {
		return nil, err
	}
	memory.Limits.Max = Pointer(max)
	return memory, nil
}

func (p *parser) ParseSignature() (*Function, error) {
	function := &Function{}
	for {
		tok, err := p.peekToken()
		if err != nil {
			return nil, err
		}
		if tok.Type != OpenParen {
			break
		}
		p.nextToken()

		name, err := p.ParseString()
		if err != nil {
			return nil, err
		}

		switch name.Capture {
		case "param":
			parameter, err := p.ParseParameter()
			if err != nil {
				return nil, err
			}
			function.Parameters = append(function.Parameters, *parameter)
		case "result":
			result, err := p.ParseResult()
			if err != nil {
				return nil, err
			}
			function.Results = append(function.Results, *result)
		default:
			return nil, p.parseError(name, fmt.Errorf("unrecognized string %s. expected 'param' or 'result'", name.Capture))
		}

		err = p.ExpectToken(CloseParen)
		if err != nil {
			return nil, err
		}
	}
	return function, nil
}

func (p *parser) ParseParameter() (*Parameter, error) {
	param := &Parameter{}
	str, err := p.ParseString()
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(str.Capture, "$") {
		param.ID = Pointer(Identifier(str.Capture))
		str, err = p.ParseString()
		if err != nil {
			return nil, err
		}
	}
	param.Type = p.ParseType(str.Capture)
	return param, nil
}

func (p *parser) ParseType(str string) Type {
	switch Type(str) {
	case I32:
		return I32
	case I64:
		return I64
	case F64:
		return F64
	case F32:
		return F32
	default:
		return ""
	}
}

func (p *parser) ParseResult() (*Result, error) {
	result := &Result{}
	str, err := p.ParseString()
	if err != nil {
		return nil, err
	}
	result.Type = p.ParseType(str.Capture)
	return result, nil
}

func (p *parser) ParseInstructions() ([]Instruction, error) {
	var instructions []Instruction
	for {
		instruction, ok, err := p.TryParseInstruction()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		instructions = append(instructions, *instruction)
	}
	return instructions, nil
}

func (p *parser) TryParseInstruction() (*Instruction, bool, error) {
	peek, err := p.peekToken()

	if err != nil {
		return nil, false, err
	}

	if peek.Type == CloseParen {
		return nil, false, nil
	}

	instruction, err := p.ParseInstruction()
	if err != nil {
		return nil, false, err
	}

	return instruction, true, nil
}

func (p *parser) ParseInstruction() (*Instruction, error) {
	str, err := p.ParseString()
	if err != nil {
		return nil, err
	}
	instruction := &Instruction{}
	split := strings.Split(str.Capture, ".")
	switch split[0] {
	case "local":
		local, err := p.ParseLocalInstruction(str.Capture)
		if err != nil {
			return nil, err
		}
		instruction.Plain = &Plain{
			Local: local,
		}
	case "i32":
		i32, err := p.ParseI32Instruction(str.Capture)
		if err != nil {
			return nil, err
		}
		instruction.Plain = &Plain{
			I32: i32,
		}
	case "i64":
	case "f32":
	case "f64":

	}

	return instruction, nil
}

func (p *parser) ParseI32Instruction(instruction string) (*I32Instruction, error) {
	split := strings.Split(instruction, ".")
	if len(split) != 2 {
		return nil, fmt.Errorf("expected i32.<operation>, found %s", instruction)
	}

	i32 := &I32Instruction{}
	operation := split[1]
	switch operation {
	case "add":
		i32.Operation = BinaryOperationAdd
	default:
		return nil, fmt.Errorf("unrecognized i32 operation %s", operation)
	}
	return i32, nil
}

func (p *parser) ParseLocalInstruction(instruction string) (*LocalInstruction, error) {
	split := strings.Split(instruction, ".")
	if len(split) != 2 {
		return nil, fmt.Errorf("expected local.<operation>, found %s", instruction)
	}

	local := &LocalInstruction{}
	operation := split[1]

	switch operation {
	case "get":
		local.Operation = LocalGet
	case "set":
		local.Operation = LocalSet
	case "tee":
		local.Operation = LocalTee
	default:
		return nil, fmt.Errorf("unrecognized local operation %s", operation)
	}

	index, err := p.ParseString()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(index.Capture, "$") {
		id := Identifier(index.Capture)
		local.ID = &id
	} else {
		i, err := strconv.Atoi(index.Capture)
		if err != nil {
			return nil, err
		}
		local.Index = &i
	}
	return local, nil
}

func (p *parser) TryParseIdentifier() (*Identifier, bool, error) {
	token, err := p.peekToken()
	if err != nil {
		return nil, false, err
	}
	if token.Type != String {
		return nil, false, nil
	}
	if strings.HasPrefix(token.Capture, "$") {
		return Pointer(Identifier(token.Capture)), true, nil
	}
	return nil, false, nil
}

func (p *parser) ParseUInt32() (uint32, error) {
	str, err := p.ParseString()
	if err != nil {
		return 0, err
	}
	result, err := strconv.ParseUint(str.Capture, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(result), err
}

func (p *parser) ParseString() (*Token, error) {
	token, err := p.nextToken()
	if err != nil {
		return nil, err
	}
	if token.Type != String {
		return nil, p.parseError(token, fmt.Errorf("expected '%s' found '%s' ", String, token.Type))
	}
	return token, nil
}

func (p *parser) ExpectString(expected string) error {
	token, err := p.ParseString()
	if err != nil {
		return err
	}
	if token.Capture != expected {
		return p.parseError(token, fmt.Errorf("expected '%s' found '%s'", expected, token.Capture))
	}
	return nil
}

func (p *parser) ExpectToken(t TokenType) error {
	token, err := p.nextToken()
	if err != nil {
		return err
	}
	if token.Type != t {
		return p.parseError(token, fmt.Errorf("expected '%s' found '%s'", t, token.Type))
	}
	return nil
}

func (p *parser) parseError(t *Token, err error) error {
	return fmt.Errorf("parse error line: %d, column: %d, position: %d, %w", t.Line+1, t.Column+1, t.Position, err)
}

func (p *parser) nextToken() (*Token, error) {
	for {
		tok, err := p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if tok.Type != Whitespace {
			return tok, nil
		}
	}
}

func (p *parser) peekToken() (*Token, error) {
	for {
		tok, err := p.lexer.Peek()
		if err != nil {
			return nil, err
		}
		if tok.Type != Whitespace {
			return tok, nil
		}
		p.lexer.Next()
	}
}
