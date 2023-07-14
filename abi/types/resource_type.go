package types

type ResourceType interface {
	Type
	resourcetype()
	DTor() DTorFunc
	Impl() *ComponentInstance
}

type DTorFunc func(uint32)

type ResourceTypeImpl struct {
	TypeImpl
	dtor DTorFunc
	impl *ComponentInstance
}

func (*ResourceTypeImpl) resourcetype() {}

func (rt *ResourceTypeImpl) DTor() DTorFunc {
	return rt.dtor
}

func (rt *ResourceTypeImpl) Impl() *ComponentInstance {
	return rt.impl
}

func NewResourceType(dtor DTorFunc, impl *ComponentInstance) ResourceType {
	return &ResourceTypeImpl{
		dtor: dtor,
		impl: impl,
	}
}
