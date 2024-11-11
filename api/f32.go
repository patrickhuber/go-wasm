package api

type F32Const float32

func (F32Const) instruction() {}

type F32Abs struct{}

func (F32Abs) instruction() {}

type F32Neg struct{}

func (F32Neg) instruction() {}

type F32Sqrt struct{}

func (F32Sqrt) instruction() {}

type F32Ceil struct{}

func (F32Ceil) instruction() {}

type F32Floor struct{}

func (F32Floor) instruction() {}

type F32Trunc struct{}

func (F32Trunc) instruction() {}

type F32Nearest struct{}

func (F32Nearest) instruction() {}

type F32Add struct{}

func (F32Add) instruction() {}

type F32Sub struct{}

func (F32Sub) instruction() {}

type F32Mul struct{}

func (F32Mul) instruction() {}

type F32Div struct{}

func (F32Div) instruction() {}

type F32Min struct{}

func (F32Min) instruction() {}

type F32Max struct{}

func (F32Max) instruction() {}

type F32CopySign struct{}

func (F32CopySign) instruction() {}

type F32Eq struct{}

func (F32Eq) instruction() {}

type F32Ne struct{}

func (F32Ne) instruction() {}

type F32Lt struct{}

func (F32Lt) instruction() {}

type F32Gt struct{}

func (F32Gt) instruction() {}

type F32Le struct{}

func (F32Le) instruction() {}

type F32Ge struct{}

func (F32Ge) instruction() {}

type F32DemoteF64 struct{}

func (F32DemoteF64) instruction() {}

type F32ConvertI32s struct{}

func (F32ConvertI32s) instruction() {}

type F32ConvertI32u struct{}

func (F32ConvertI32u) instruction() {}

type F32ReinterpretI32 struct{}
