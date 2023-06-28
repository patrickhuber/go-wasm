package types

import (
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/internal/collections/stack"
)

type HandleElem struct {
	Rep       int
	LendCount int
	Context   *CallContext // always null for OwnHandle
}

type HandleTable struct {
	Array []*HandleElem
	Free  []uint32
}

func (ht *HandleTable) Add(handle *HandleElem, t ValType) uint32 {
	free, i, ok := stack.Pop(ht.Free)
	ht.Free = free
	if ok {
		ht.Array[i] = handle
	} else {
		ht.Free = stack.Push(ht.Free, uint32(len(ht.Free)))
		ht.Array = append(ht.Array, handle)
	}
	switch t.Kind() {
	case kind.Borrow:
		handle.Context.BorrowCount += 1
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

func (ht *HandleTable) Remove(i uint32, t ValType, drop bool) (*HandleElem, error) {
	// null dereference?
	h, err := ht.Get(i)
	if err != nil {
		return nil, err
	}

	// open handles?
	if h.LendCount != 0 {
		return nil, TrapWith("handle table end count != 0")
	}

	switch t.Kind() {
	case kind.Own:
		own, ok := t.(*Own)
		if !ok {
			return nil, TrapWith("unable to cast %T to types.Own", t)
		}
		if !drop || own.ResourceType.DTor == nil {
			break
		}
		if !own.ResourceType.Impl.MayEnter {
			return nil, TrapWith("MayEnter is false")
		}
		(*own.ResourceType.DTor)(h.Rep)

	case kind.Borrow:
		_, ok := t.(*Borrow)
		if !ok {
			return nil, TrapWith("unable to cast %T to types.Borrow", t)
		}
		h.Context.BorrowCount -= 1
	}

	ht.Array[i] = nil
	ht.Free = stack.Push(ht.Free, uint32(len(ht.Free)))
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

func resourceType(t ValType) (*ResourceType, error) {

	switch t.Kind() {
	case kind.Borrow:
		b, ok := t.(*Borrow)
		if !ok {
			return nil, TrapWith("unable to cast %T to types.Borrow", t)
		}
		return b.ResourceType, nil
	case kind.Own:
		o, ok := t.(*Own)
		if !ok {
			return nil, TrapWith("unable to cast %T to types.Own", t)
		}
		return o.ResourceType, nil
	}
	return nil, TrapWith("unrecognized kind %s", t.Kind())
}

func (ht *HandleTables) Add(handle *HandleElem, t ValType) (uint32, error) {
	resourceType, err := resourceType(t)
	if err != nil {
		return 0, err
	}
	return ht.Table(*resourceType).Add(handle, t), nil
}

func (ht *HandleTables) Get(i uint32, resourceType *ResourceType) (*HandleElem, error) {
	return ht.Table(*resourceType).Get(i)
}

func (ht *HandleTables) Remove(i uint32, t ValType) (*HandleElem, error) {
	resourceType, err := resourceType(t)
	if err != nil {
		return nil, err
	}
	return ht.Table(*resourceType).Remove(i, t, false)
}
