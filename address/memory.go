package address

type Memory struct {
	Address uint32
}

func (*Memory) address()       {}
func (*Memory) externalvalue() {}
