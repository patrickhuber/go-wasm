package types

type CallContext struct {
	Options     *CanonicalOptions
	Instance    *ComponentInstance
	Lenders     []*HandleElem
	BorrowCount int
}

func (cx *CallContext) LiftBorrowFrom(lendingHandle *HandleElem) {
	lendingHandle.LendCount += 1
	cx.Lenders = append(cx.Lenders, lendingHandle)
}

func (cx *CallContext) RemoveBorrowFromTable() {
	cx.BorrowCount -= 1
}

func (cx *CallContext) AddBorrowCountToTable() {
	cx.BorrowCount += 1
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
