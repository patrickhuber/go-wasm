package types

import (
	"github.com/patrickhuber/go-wasm/abi/kind"
	"github.com/patrickhuber/go-wasm/collections/stack"
)

type Handle struct {
	Rep       int
	LendCount int
	Context   *Context // always null for OwnHandle
}

type HandleTable struct {
	Array []*Handle
	Free  []int
}

func (ht *HandleTable) Add(handle *Handle, t ValType) int {
	free, i, ok := stack.Pop(ht.Free)
	ht.Free = free
	if ok {
		ht.Array[i] = handle
	} else {
		ht.Free = stack.Push(ht.Free, len(ht.Free))
		ht.Array = append(ht.Array, handle)
	}
	switch t.Kind() {
	case kind.Borrow:
		handle.Context.BorrowCount += 1
	}
	return i
}

func (ht *HandleTable) Get(i int) (*Handle, error) {
	if err := TrapIf(i >= len(ht.Array)); err != nil {
		return nil, err
	}
	handle := ht.Array[i]
	if err := TrapIf(handle == nil); err != nil {
		return nil, err
	}
	return handle, nil
}

func (ht *HandleTable) TransferOrDrop(i int, t ValType, drop bool) (*Handle, error) {
	// null dereference?
	h, err := ht.Get(i)
	if err != nil {
		return nil, err
	}

	// open handles?
	err = TrapIf(h.LendCount != 0)
	if err != nil {
		return nil, err
	}

	switch t.Kind() {
	case kind.Own:
		own, ok := t.(*Own)
		err = TrapIf(!ok)
		if err != nil {
			return nil, err
		}
		if !drop || own.ResourceType.DTor == nil {
			break
		}
		err = TrapIf(!own.ResourceType.Impl.MayEnter)
		if err != nil {
			return nil, err
		}
		(*own.ResourceType.DTor)(h.Rep)

	case kind.Borrow:
		_, ok := t.(*Borrow)
		err = TrapIf(!ok)
		if err != nil {
			return nil, err
		}
		h.Context.BorrowCount -= 1
	}

	ht.Array[i] = nil
	ht.Free = stack.Push(ht.Free, len(ht.Free))
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
		err := TrapIf(!ok)
		if err != nil {
			return nil, err
		}
		return b.ResourceType, nil
	case kind.Own:
		o, ok := t.(*Own)
		err := TrapIf(!ok)
		if err != nil {
			return nil, err
		}
		return o.ResourceType, nil
	}
	return nil, &Trap{}
}

func (ht *HandleTables) Add(handle *Handle, t ValType) (int, error) {
	resourceType, err := resourceType(t)
	if err != nil {
		return 0, err
	}
	return ht.Table(*resourceType).Add(handle, t), nil
}

func (ht *HandleTables) Get(i int, resourceType *ResourceType) (*Handle, error) {
	return ht.Table(*resourceType).Get(i)
}

func (ht *HandleTables) Transfer(i int, t ValType) (*Handle, error) {
	resourceType, err := resourceType(t)
	if err != nil {
		return nil, err
	}
	return ht.Table(*resourceType).TransferOrDrop(i, t, false)
}

func (ht *HandleTables) Drop(i int, t ValType) error {
	resourceType, err := resourceType(t)
	if err != nil {
		return err
	}
	_, err = ht.Table(*resourceType).TransferOrDrop(i, t, true)
	return err
}
