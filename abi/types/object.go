package types

type Object interface {
	object()
}

type ObjectImpl struct {
}

func (*ObjectImpl) object() {
}
