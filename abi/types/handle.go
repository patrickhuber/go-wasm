package types

import (
	"github.com/patrickhuber/go-wasm/internal/collections/stack"
)

type HandleElem struct {
	Rep       int
	LendCount int
	Own       bool
	Scope     *CallContext // always null for OwnHandle
}

type HandleTable struct {
	Array []*HandleElem
	Free  []uint32
}

func (ht *HandleTable) Add(handle *HandleElem) uint32 {
	free, i, ok := stack.Pop(ht.Free)
	ht.Free = free
	if ok {
		ht.Array[i] = handle
	} else {
		ht.Free = stack.Push(ht.Free, uint32(len(ht.Free)))
		ht.Array = append(ht.Array, handle)
	}
	return i
}

func (ht *HandleTable) Get(i uint32) (*HandleElem, error) {
	if i >= uint32(len(ht.Array)) {
		return nil, TrapWith("index is greater than handle table length")
	}
	handle := ht.Array[i]
	if handle == nil {
		return nil, TrapWith("handle %d is nil", i)
	}
	return handle, nil
}

func (ht *HandleTable) Remove(rt ResourceType, i uint32) (*HandleElem, error) {
	// null dereference?
	h, err := ht.Get(i)
	if err != nil {
		return nil, err
	}

	// open handles?
	if h.LendCount != 0 {
		return nil, TrapWith("handle table end count != 0")
	}
	ht.Array[i] = nil
	ht.Free = stack.Push(ht.Free, uint32(len(ht.Free)))
	if h.Scope != nil {
		h.Scope.RemoveBorrowFromTable()
	}
	return h, nil
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

func (ht *HandleTables) Add(handle *HandleElem, rt ResourceType) (uint32, error) {
	return ht.Table(rt).Add(handle), nil
}

func (ht *HandleTables) Get(rt ResourceType, i uint32) (*HandleElem, error) {
	return ht.Table(rt).Get(i)
}

func (ht *HandleTables) Remove(rt ResourceType, i uint32) (*HandleElem, error) {
	return ht.Table(rt).Remove(rt, i)
}
