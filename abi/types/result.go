package types

type Result struct {
	OK    ValType
	Error ValType
}

func (r *Result) Size() uint32 {
	return r.Despecialize().Size()
}

func (r *Result) Alignment() uint32 {
	return r.Despecialize().Alignment()
}

func (r *Result) Despecialize() ValType {
	cases := []Case{
		{
			Label: "ok",
			Type:  r.OK,
		},
		{
			Label: "error",
			Type:  r.Error,
		},
	}
	return &Variant{
		Cases: cases,
	}
}
