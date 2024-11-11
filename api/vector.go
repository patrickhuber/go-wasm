package api

type Vector interface {
	vector()
}

type Vector128 struct{}

func (Vector128) vector() {}
