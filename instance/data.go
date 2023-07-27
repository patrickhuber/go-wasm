package instance

type Data struct {
	Data []byte
}

func (*Data) instance() {}
