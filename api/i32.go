package api

type I32Const uint32

func (I32Const) instruction() {}

type I32Add struct{}

func (I32Add) instruction() {}

type I32Sub struct{}

func (I32Sub) instruction() {}

type I32Mul struct{}

func (I32Mul) instruction() {}

type I32Div struct{}

func (I32Div) instruction() {}

type U32Div struct{}

func (U32Div) instruction() {}

type I32Rem struct{}

func (I32Rem) instruction() {}

type U32Rem struct{}

func (U32Rem) instruction() {}

type I32And struct{}

func (I32And) instruction() {}

type I32Or struct{}

func (I32Or) instruction() {}

type I32Xor struct{}

func (I32Xor) instruction() {}

type I32Shl struct{}

func (I32Shl) instruction() {}

type I32Shr struct{}

func (I32Shr) instruction() {}

type U32Shr struct{}

func (U32Shr) instruction() {}

type I32Rotl struct{}

func (I32Rotl) instruction() {}

type I32Rotr struct{}

func (I32Rotr) instruction() {}

type I32Eqz struct{}

func (I32Eqz) instruction() {}

type I32Eq struct{}

func (I32Eq) instruction() {}

type I32Ne struct{}

func (I32Ne) instruction() {}

type I32Lt struct{}

func (I32Lt) instruction() {}

type U32Lt struct{}

func (U32Lt) instruction() {}

type I32Gt struct{}

func (I32Gt) instruction() {}

type U32Gt struct{}

func (U32Gt) instruction() {}

type I32Le struct{}

func (I32Le) instruction() {}

type U32Le struct{}

func (U32Le) instruction() {}

type I32Ge struct{}

func (I32Ge) instruction() {}

type U32Ge struct{}

func (U32Ge) instruction() {}
