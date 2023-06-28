package types

type CallContext struct {
	Options     *CanonicalOptions
	Instance    *ComponentInstance
	Lenders     []*HandleElem
	BorrowCount int
}
