package binary

var Magic = []byte{0x00, 0x61, 0x73, 0x6d}

const ModuleVersion uint16 = 0x01
const ComponentVersion uint16 = 0x0a

type SectionID uint8

const (
	CustomSectionID   SectionID = 0
	TypeSectionID     SectionID = 1
	FunctionSectionID SectionID = 3
	ExportSectionID   SectionID = 7
	CodeSectionID     SectionID = 10
)

type Section struct {
	ID   SectionID
	Data []byte
}

type ValType byte

const I32 ValType = 0x7f
const I64 ValType = 0x7e
const F32 ValType = 0x7d
const F64 ValType = 0x7c

type ExportKind byte

const FuncExportKind ExportKind = 0x00
const TableExportKind ExportKind = 0x01
const MemoryExportKind ExportKind = 0x02
const GlobalExportKind ExportKind = 0x03
