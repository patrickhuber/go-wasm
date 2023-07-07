package indicies

type Index interface {
	index()
}

type Function uint32

func (*Function) index() {}

type Type uint32

func (*Type) index() {}

type Table uint32

func (*Table) index() {}

type Memory uint32

func (*Memory) index() {}

type Global uint32

func (*Global) index() {}

type Element uint32

func (*Element) index() {}

type Data uint32

func (*Data) index() {}

type Local uint32

func (*Local) index() {}

type Label uint32

func (*Label) index() {}
