package types

type Borrow interface {
	ValType
	ResourceType() ResourceType
	borrow()
}

type BorrowImpl struct {
	ValTypeImpl
	resourceType ResourceType
}

func (*BorrowImpl) borrow() {}

func (b *BorrowImpl) ResourceType() ResourceType {
	return b.resourceType
}

func NewBorrow(rt ResourceType) Borrow {
	return &BorrowImpl{
		resourceType: rt,
	}
}
