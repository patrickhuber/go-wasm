package address

type Function struct {
	Address uint32
}

func (*Function) address()       {}
func (*Function) externalvalue() {}
