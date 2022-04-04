package wat

type ModuleBuilder interface {
	Build() *Module
}

type SectionBuilder interface {
	Function(func(FunctionBuilder))
	Memory(func(MemoryBuilder))
}

type MemoryBuilder interface {
	Limits(min uint32)
}

type FunctionBuilder interface {
	ID(id string)
	Parameters(func(ParametersBuilder))
	Results(func(ResultsBuilder))
	Instructions(func(InstructionsBuilder))
}

type ParametersBuilder interface {
	Parameter(t Type) ParameterBuilder
}

type ParameterBuilder interface {
	ID(id string)
}

type ResultsBuilder interface {
	Result(t Type)
}

type InstructionsBuilder interface {
	Local(op LocalOperation) LocalInstructionBuilder
	I32Add()
}

type LocalInstructionBuilder interface {
	ID(name string) LocalInstructionBuilder
	Index(index int) LocalInstructionBuilder
}

func NewModule(sections func(s SectionBuilder)) ModuleBuilder {
	m := &Module{}
	sectionBuilder := &sectionBuilder{}
	sections(sectionBuilder)
	for _, s := range sectionBuilder.sections {
		if s.Function != nil {
			m.Functions = append(m.Functions, s)
		}
		if s.Memory != nil {
			m.Memory = append(m.Memory, s)
		}
	}
	return &moduleBuilder{
		module: m,
	}
}

type moduleBuilder struct {
	module *Module
}

func (b *moduleBuilder) Build() *Module {
	return b.module
}

type sectionBuilder struct {
	sections []Section
}

func (b *sectionBuilder) Function(f func(FunctionBuilder)) {
	function := &functionBuilder{
		function: &Function{},
	}
	f(function)
	b.sections = append(b.sections, Section{
		Function: function.function,
	})
}

func (b *sectionBuilder) Memory(m func(MemoryBuilder)) {
	memory := &memoryBuilder{
		memory: &Memory{},
	}
	m(memory)
	b.sections = append(b.sections, Section{
		Memory: memory.memory,
	})
}

type functionBuilder struct {
	function *Function
}

func (b *functionBuilder) ID(id string) {
	identifier := Identifier(id)
	b.function.ID = &identifier
}

func (b *functionBuilder) Parameters(p func(ParametersBuilder)) {
	parameters := &parametersBuilder{}
	p(parameters)
	for _, param := range parameters.parameters {
		b.function.Parameters = append(b.function.Parameters, *param)
	}
}

func (b *functionBuilder) Results(r func(ResultsBuilder)) {
	results := &results{}
	r(results)
	for _, result := range results.results {
		b.function.Results = append(b.function.Results, result)
	}
}

func (b *functionBuilder) Instructions(i func(InstructionsBuilder)) {
	instructions := &instructionsBuilder{}
	i(instructions)
	for _, instruction := range instructions.instructions {
		b.function.Instructions = append(b.function.Instructions, instruction)
	}
}

type memoryBuilder struct {
	memory *Memory
}

func (b *memoryBuilder) Limits(min uint32) {
	b.memory.Limits = Limits{
		Min: min,
	}
}

type parametersBuilder struct {
	parameters []*Parameter
}

func (b *parametersBuilder) Parameter(t Type) ParameterBuilder {
	p := &parameterBuilder{
		parameter: &Parameter{
			Type: t,
		},
	}
	b.parameters = append(b.parameters, p.parameter)
	return p
}

type parameterBuilder struct {
	parameter *Parameter
}

func (b *parameterBuilder) ID(id string) {
	identifer := Identifier(id)
	b.parameter.ID = &identifer
}

type results struct {
	results []Result
}

func (b *results) Result(t Type) {
	result := Result{
		Type: t,
	}
	b.results = append(b.results, result)
}

type instructionsBuilder struct {
	instructions []Instruction
}

func (b *instructionsBuilder) Local(op LocalOperation) LocalInstructionBuilder {
	inst := &localInstructionBuilder{
		local: &LocalInstruction{
			Operation: op,
		},
	}
	b.instructions = append(b.instructions, Instruction{
		Plain: &Plain{
			Local: inst.local,
		},
	})
	return inst
}

func (b *instructionsBuilder) I32Add() {
	b.instructions = append(b.instructions, Instruction{
		Plain: &Plain{
			I32: &I32Instruction{
				Operation: BinaryOperationAdd,
			},
		},
	})
}

type localInstructionBuilder struct {
	local *LocalInstruction
}

func (b *localInstructionBuilder) ID(id string) LocalInstructionBuilder {
	identifier := Identifier(id)
	b.local.ID = &identifier
	return b
}

func (b *localInstructionBuilder) Index(index int) LocalInstructionBuilder {
	b.local.Index = &index
	return b
}
