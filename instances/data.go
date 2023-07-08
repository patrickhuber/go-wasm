package instances

type Data struct {
	Data []byte
}

func (*Data) instance() {}
