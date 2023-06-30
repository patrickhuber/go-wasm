package types

type Own interface {
	ValType
	ResourceType() ResourceType
	own()
}

type OwnImpl struct {
	ValTypeImpl
	resourceType ResourceType
}

func (*OwnImpl) own() {}

func (o *OwnImpl) ResourceType() ResourceType {
	return o.resourceType
}

func NewOwn(rt ResourceType) Own {
	return &OwnImpl{
		resourceType: rt,
	}
}
