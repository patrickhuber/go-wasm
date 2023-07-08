package address

type Global struct {
	Address uint32
}

func (*Global) address()       {}
func (*Global) externalvalue() {}
