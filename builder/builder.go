package builder

import "github.com/patrickhuber/go-wasm"

type Module interface {
	Build() *wasm.Module
}

type Section interface {
	Function(func(Function))
	Memory(func(Memory))
}

type Memory interface {
	Limits(min uint32)
}

type Function interface {
	Parameters(func(Parameters))
	Results(func(Results))
	Instructions(func(Instructions))
}

type Parameters interface {
	Parameter(t wasm.Type) Parameter
}

type Parameter interface {
	ID(id string)
}

type Results interface {
	Result(t wasm.Type)
}

type Instructions interface {
	Local(op wasm.LocalOperation) LocalInstruction
	I32Add()
}

type LocalInstruction interface {
	ID(name string) LocalInstruction
	Index(index int) LocalInstruction
}

func NewModule(sections func(s Section)) Module {
	m := &wasm.Module{}
	sectionBuilder := &sectionBuilder{}
	sections(sectionBuilder)
	for _, s := range sectionBuilder.sections {
		if s.Function != nil {
			m.Functions = append(m.Functions, *s.Function)
		}
		if s.Memory != nil {
			m.Memory = append(m.Memory, *s.Memory)
		}
	}
	return &module{
		module: m,
	}
}

type module struct {
	module *wasm.Module
}

func (b *module) Build() *wasm.Module {
	return b.module
}

type sectionBuilder struct {
	sections []wasm.Section
}

func (b *sectionBuilder) Function(f func(Function)) {
	function := &function{
		function: &wasm.Function{},
	}
	f(function)
	b.sections = append(b.sections, wasm.Section{
		Function: function.function,
	})
}

func (b *sectionBuilder) Memory(m func(Memory)) {
	memory := &memory{
		memory: &wasm.Memory{},
	}
	m(memory)
	b.sections = append(b.sections, wasm.Section{
		Memory: memory.memory,
	})
}

type function struct {
	function *wasm.Function
}

func (b *function) Parameters(p func(Parameters)) {
	parameters := &parameters{}
	p(parameters)
	for _, param := range parameters.parameters {
		b.function.Parameters = append(b.function.Parameters, *param)
	}
}

func (b *function) Results(r func(Results)) {
	results := &results{}
	r(results)
	for _, result := range results.results {
		b.function.Results = append(b.function.Results, result)
	}
}

func (b *function) Instructions(i func(Instructions)) {
	instructions := &instructions{}
	i(instructions)
	for _, instruction := range instructions.instructions {
		b.function.Instructions = append(b.function.Instructions, instruction)
	}
}

type memory struct {
	memory *wasm.Memory
}

func (b *memory) Limits(min uint32) {
	b.memory.Limits = wasm.Limits{
		Min: min,
	}
}

type parameters struct {
	parameters []*wasm.Parameter
}

func (b *parameters) Parameter(t wasm.Type) Parameter {
	p := &parameter{
		parameter: &wasm.Parameter{
			Type: t,
		},
	}
	b.parameters = append(b.parameters, p.parameter)
	return p
}

type parameter struct {
	parameter *wasm.Parameter
}

func (b *parameter) ID(id string) {
	identifer := wasm.Identifier(id)
	b.parameter.ID = &identifer
}

type results struct {
	results []wasm.Result
}

func (b *results) Result(t wasm.Type) {
	result := wasm.Result{
		Type: t,
	}
	b.results = append(b.results, result)
}

type instructions struct {
	instructions []wasm.Instruction
}

func (b *instructions) Local(op wasm.LocalOperation) LocalInstruction {
	inst := &localInstruction{
		local: &wasm.LocalInstruction{
			Operation: op,
		},
	}
	b.instructions = append(b.instructions, wasm.Instruction{
		Plain: &wasm.Plain{
			Local: inst.local,
		},
	})
	return inst
}

type localInstruction struct {
	local *wasm.LocalInstruction
}

func (b *localInstruction) ID(id string) LocalInstruction {
	identifier := wasm.Identifier(id)
	b.local.ID = &identifier
	return b
}

func (b *localInstruction) Index(index int) LocalInstruction {
	b.local.Index = &index
	return b
}

func (b *instructions) I32Add() {
	b.instructions = append(b.instructions, wasm.Instruction{
		Plain: &wasm.Plain{
			I32: &wasm.I32Instruction{
				Operation: wasm.BinaryOperationAdd,
			},
		},
	})
}
