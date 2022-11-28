package wasm

type Version uint16
type Layer uint16

const (
	Magic               uint32  = 0x0061_736d // big endian
	Version1            Version = 1
	VersionExperimental Version = 10
	LayerCore           Layer   = 0
	LayerComponent      Layer   = 1
)

type Header struct {
	Magic   uint32
	Version Version
	Layer   Layer
}

func NewModuleHeader() *Header {
	return &Header{
		Magic:   Magic,
		Version: Version1,
		Layer:   LayerCore,
	}
}

func NewComponentHeader() *Header {
	return &Header{
		Magic:   Magic,
		Version: VersionExperimental,
		Layer:   LayerComponent,
	}
}
