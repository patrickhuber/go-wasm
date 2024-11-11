package api

type Number interface {
	number()
}

type I32 struct{}

func (I32) number() {}
func (I32) value()  {}

type I64 struct{}

func (I64) number() {}
func (I64) value()  {}

type F32 struct{}

func (F32) number() {}
func (F32) value()  {}

type F64 struct{}

func (F64) number() {}
func (F64) value()  {}
