package types

type Handle struct {
	Rep       int
	LendCount int
}

type OwnHandle struct {
	Handle
}

type BorrowHandle struct {
	Handle
	Context *Context
}

type HandleTable struct {
	Array []Handle
	Free  []int
}

func (ht *HandleTable) Add(h Handle, t ValType) int {
	return 0
}

func (ht *HandleTable) Get(i int) (Handle, error) {
	return Handle{}, nil
}

func (ht *HandleTable) TransferOrDrop(i int, t ValType, drop any) {

}

type HandleTables struct {
	ResourceTypeToTable map[ResourceType]*HandleTable
}

func (ht *HandleTables) Table(rt ResourceType) *HandleTable {
	t, ok := ht.ResourceTypeToTable[rt]
	if !ok {
		t = &HandleTable{}
		ht.ResourceTypeToTable[rt] = t
	}
	return t
}
