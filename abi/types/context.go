package types

type CallContext struct {
	Options     *CanonicalOptions
	Instance    *ComponentInstance
	Lenders     []*HandleElem
	BorrowCount int
}

func (cx *CallContext) ExitCall() error {
	if cx.BorrowCount != 0 {
		return TrapWith("borrow count != 0")
	}
	for _, h := range cx.Lenders {
		h.LendCount -= 1
	}
	return nil
}
