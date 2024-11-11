package api

type I64Const uint64

func (I64Const) instruction() {}

type I64Add struct{}

func (I64Add) instruction() {}

type I64Sub struct{}

func (I64Sub) instruction() {}

type I64Mul struct{}

func (I64Mul) instruction() {}

type I64Div struct{}

func (I64Div) instruction() {}

type U64Div struct{}

func (U64Div) instruction() {}

type I64Rem struct{}

func (I64Rem) instruction() {}

type U64Rem struct{}

func (U64Rem) instruction() {}

type I64And struct{}

func (I64And) instruction() {}

type I64Or struct{}

func (I64Or) instruction() {}

type I64Xor struct{}

func (I64Xor) instruction() {}

type I64Shl struct{}

func (I64Shl) instruction() {}

type I64Shr struct{}

func (I64Shr) instruction() {}

type U64Shr struct{}

func (U64Shr) instruction() {}

type I64Rotl struct{}

func (I64Rotl) instruction() {}

type I64Rotr struct{}

func (I64Rotr) instruction() {}

type I64Eqz struct{}

func (I64Eqz) instruction() {}

type I64Eq struct{}

func (I64Eq) instruction() {}

type I64Ne struct{}

func (I64Ne) instruction() {}

type I64Lt struct{}

func (I64Lt) instruction() {}

type U64Lt struct{}

func (U64Lt) instruction() {}

type I64Gt struct{}

func (I64Gt) instruction() {}

type U64Gt struct{}

func (U64Gt) instruction() {}

type I64Le struct{}

func (I64Le) instruction() {}

type U64Le struct{}

func (U64Le) instruction() {}

type I64Ge struct{}

func (I64Ge) instruction() {}

type U64Ge struct{}

func (U64Ge) instruction() {}
