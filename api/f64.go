package api

type F64Const float64

func (F64Const) instruction() {}

type F64Abs struct{}

func (F64Abs) instruction() {}

type F64Neg struct{}

func (F64Neg) instruction() {}

type F64Sqrt struct{}

func (F64Sqrt) instruction() {}

type F64Ceil struct{}

func (F64Ceil) instruction() {}

type F64Floor struct{}

func (F64Floor) instruction() {}

type F64Trunc struct{}

func (F64Trunc) instruction() {}

type F64Nearest struct{}

func (F64Nearest) instruction() {}

type F64Add struct{}

func (F64Add) instruction() {}

type F64Sub struct{}

func (F64Sub) instruction() {}

type F64Mul struct{}

func (F64Mul) instruction() {}

type F64Div struct{}

func (F64Div) instruction() {}

type F64Min struct{}

func (F64Min) instruction() {}

type F64Max struct{}

func (F64Max) instruction() {}

type F64CopySign struct{}

func (F64CopySign) instruction() {}

type F64Eq struct{}

func (F64Eq) instruction() {}

type F64Ne struct{}

func (F64Ne) instruction() {}

type F64Lt struct{}

func (F64Lt) instruction() {}

type F64Gt struct{}

func (F64Gt) instruction() {}

type F64Le struct{}

func (F64Le) instruction() {}

type F64Ge struct{}

func (F64Ge) instruction() {}

type F64ConvertI32s struct{}

func (F64ConvertI32s) instruction() {}

type F64PromoteF32 struct{}

func (F64PromoteF32) instruction() {}

type F64ConvertI32u struct{}

func (F64ConvertI32u) instruction() {}

type F64ConvertI64u struct{}

func (F64ConvertI64u) instruction() {}

type F64ReinterpretI32 struct{}

func (F64ReinterpretI32) instruction() {}

type F64ReinterpretI64 struct{}

func (F64ReinterpretI64) instruction() {}
