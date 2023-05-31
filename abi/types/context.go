package types

type Context struct {
	Options     *CanonicalOptions
	Instance    *ComponentInstance
	Lenders     []Handle
	BorrowCount int
}
