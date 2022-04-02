package builder

import (
	"github.com/patrickhuber/go-wasm/model"
)

type Module interface {
	Build() *model.Module
}

type Section interface {
	Function(func(Function))
	Memory(func(Memory))
}

type Memory interface {
	Limits(min uint32)
}

type Function interface {
	ID(id string)
	Parameters(func(Parameters))
	Results(func(Results))
	Instructions(func(Instructions))
}

type Parameters interface {
	Parameter(t model.Type) Parameter
}

type Parameter interface {
	ID(id string)
}

type Results interface {
	Result(t model.Type)
}

type Instructions interface {
	Local(op model.LocalOperation) LocalInstruction
	I32Add()
}

type LocalInstruction interface {
	ID(name string) LocalInstruction
	Index(index int) LocalInstruction
}

func NewModule(sections func(s Section)) Module {
	m := &model.Module{}
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
	module *model.Module
}

func (b *module) Build() *model.Module {
	return b.module
}

type sectionBuilder struct {
	sections []model.Section
}

func (b *sectionBuilder) Function(f func(Function)) {
	function := &function{
		function: &model.Function{},
	}
	f(function)
	b.sections = append(b.sections, model.Section{
		Function: function.function,
	})
}

func (b *sectionBuilder) Memory(m func(Memory)) {
	memory := &memory{
		memory: &model.Memory{},
	}
	m(memory)
	b.sections = append(b.sections, model.Section{
		Memory: memory.memory,
	})
}

type function struct {
	function *model.Function
}

func (b *function) ID(id string) {
	identifier := model.Identifier(id)
	b.function.ID = &identifier
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
	memory *model.Memory
}

func (b *memory) Limits(min uint32) {
	b.memory.Limits = model.Limits{
		Min: min,
	}
}

type parameters struct {
	parameters []*model.Parameter
}

func (b *parameters) Parameter(t model.Type) Parameter {
	p := &parameter{
		parameter: &model.Parameter{
			Type: t,
		},
	}
	b.parameters = append(b.parameters, p.parameter)
	return p
}

type parameter struct {
	parameter *model.Parameter
}

func (b *parameter) ID(id string) {
	identifer := model.Identifier(id)
	b.parameter.ID = &identifer
}

type results struct {
	results []model.Result
}

func (b *results) Result(t model.Type) {
	result := model.Result{
		Type: t,
	}
	b.results = append(b.results, result)
}

type instructions struct {
	instructions []model.Instruction
}

func (b *instructions) Local(op model.LocalOperation) LocalInstruction {
	inst := &localInstruction{
		local: &model.LocalInstruction{
			Operation: op,
		},
	}
	b.instructions = append(b.instructions, model.Instruction{
		Plain: &model.Plain{
			Local: inst.local,
		},
	})
	return inst
}

type localInstruction struct {
	local *model.LocalInstruction
}

func (b *localInstruction) ID(id string) LocalInstruction {
	identifier := model.Identifier(id)
	b.local.ID = &identifier
	return b
}

func (b *localInstruction) Index(index int) LocalInstruction {
	b.local.Index = &index
	return b
}

func (b *instructions) I32Add() {
	b.instructions = append(b.instructions, model.Instruction{
		Plain: &model.Plain{
			I32: &model.I32Instruction{
				Operation: model.BinaryOperationAdd,
			},
		},
	})
}
