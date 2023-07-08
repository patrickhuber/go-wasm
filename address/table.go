package address

type Table struct {
	Address uint32
}

func (*Table) address()       {}
func (*Table) externalvalue() {}
