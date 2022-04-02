package wasm

const (
	Magic   uint32 = 0x0061_736d // big endian
	Version uint32 = 0x0000_0001
)

type Module struct {
	Magic    uint32
	Version  uint32
	Sections []Section
}
