package binary

type Header struct{}
type Document struct {
	Header *Header
	Root   Root
}

type Root interface {
	root()
}
type Component struct{}

func (Component) root() {}

type Module struct {
	Sections []Section
}

func (Module) root() {}

type Section interface {
	section()
}

type Function struct{}

func (Function) section() {}
