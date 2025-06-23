package runtime

type ExternVal interface {
	externVal()
}

type FuncExternVal struct{}
