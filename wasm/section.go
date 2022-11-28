package wasm

type SectionType byte

const (
	CustomSectionType SectionType = iota
	TypeSectionType
	ImportSectionType
	FuncSectionType
	TableSectionType
	MemSectionType
	GlobalSectionType
	ExportSectionType
	StartSectionType
	ElemSectionType
	CodeSectionType
	DataSectionType
)

type Section struct {
	ID       SectionType
	Size     uint32
	Custom   *CustomSection
	Start    *StartSection
	Type     *TypeSection
	Import   *ImportSection
	Function *FunctionSection
	Code     *CodeSection
	Export   *ExportSection
}

type CustomSection struct {
	Name *NameSection
}

type SubSectionID byte

const (
	SubSectionName          SubSectionID = 0
	SubSectionFunctionNames SubSectionID = 1
	SubSectionLocalNames    SubSectionID = 2
)

type NameSection struct {
	ID   SubSectionID
	Key  string
	Name *string
}

type StartSection struct {
	Size uint32
	// Function index to the start function
	Function uint32
}

type TypeSection struct {
	Types []Type
}

type Type struct {
	Parameters *ResultType
	Results    *ResultType
}

type ResultType struct {
	Values []*ValueType
}

type ValueType struct {
	NumberType    *NumberType
	ReferenceType *ReferenceType
}

type NumberType byte

const (
	I32 NumberType = 0x7f
	I64 NumberType = 0x7e
	F32 NumberType = 0x7d
	F64 NumberType = 0x7c
)

type ReferenceType byte

const (
	FuncRef   ReferenceType = 0x70
	ExternRef ReferenceType = 0x6f
)

type ImportSection struct{}

type FunctionSection struct {
	Types []uint32
}

type CodeSection struct {
	Codes []Code
}

type Code struct {
	Size       uint32
	Locals     []Locals
	Expression []Instruction
}

type Locals struct {
	Count uint32
	Type  *ValueType
}

type ExportSection struct{}
type Export struct{}
