package address

type External struct {
	Address uint32
}

func (*External) address() {}
