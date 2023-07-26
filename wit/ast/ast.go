package ast

import (
	"github.com/patrickhuber/go-types"
)

type Ast struct {
	PackageName types.Option[PackageName]
	Items       []AstItem
}

type PackageName struct {
	Namespace []rune
	Name      []rune
	Version   types.Option[Version]
}

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Pre   []rune
	Build []rune
}

type AstItem struct {
	Interface *Interface
	World     *World
	Use       *TopLevelUse
}

type Interface struct {
	Name  []rune
	Items []InterfaceItem
}

type InterfaceItem struct {
	TypeDef *TypeDef
	Func    *NamedFunc
	Use     *Use
}

type NamedFunc struct {
	Name []rune
	Func *Func
}

type Func struct {
	Params  []Parameter
	Results *ResultList
}

type ResultList struct {
	Named     []Parameter
	Anonymous Type
}

type Parameter struct {
	Id   []rune
	Type Type
}

type World struct {
	Id    []rune
	Items []WorldItem
}

type WorldItem interface {
	worldItem()
}

type Export interface {
	worldItem()
	exp()
}

type ExportExternType struct {
	Id         []rune
	ExternType *ExternType
}

func (imp *ExportExternType) imp()       {}
func (imp *ExportExternType) worldItem() {}

type Import interface {
	worldItem()
	imp()
}

type ImportInterface struct {
	Interface *Interface
}

func (imp *ImportInterface) imp()       {}
func (imp *ImportInterface) worldItem() {}

type ImportExternType struct {
	Id         []rune
	ExternType *ExternType
}

func (imp *ImportExternType) imp()       {}
func (imp *ImportExternType) worldItem() {}

type ExternType struct {
	Func      *Func
	Interface *Interface
	UsePath   *UsePath
}

type Use struct {
	From  *UsePath
	Names []UseName
}

type UsePath struct {
	Id      []rune
	Package struct {
		Id   *PackageName
		Name []rune
	}
}

type UseName struct {
	Name []rune
	As   types.Option[[]rune]
}

type TypeDef struct {
	Name []rune
	Type Type
}

type Include struct{}

type TopLevelUse struct{}

type Type interface {
	ty()
}

type TypeImpl struct{}

func (t *TypeImpl) ty() {}

type U8 struct{ TypeImpl }
type U16 struct{ TypeImpl }
type U32 struct{ TypeImpl }
type U64 struct{ TypeImpl }
type S8 struct{ TypeImpl }
type S16 struct{ TypeImpl }
type S32 struct{ TypeImpl }
type S64 struct{ TypeImpl }
type Float32 struct{ TypeImpl }
type Float64 struct{ TypeImpl }
type Char struct{ TypeImpl }
type Bool struct{ TypeImpl }
type String struct{ TypeImpl }

type Tuple struct {
	TypeImpl
	Types []Type
}

type List struct {
	TypeImpl
	Type Type
}

type Option struct {
	TypeImpl
	Type Type
}
type Result struct {
	TypeImpl
	Ok    types.Option[Type]
	Error types.Option[Type]
}

type Handle struct{ TypeImpl }

type Id struct {
	TypeImpl
	Value []rune
}

type Stream struct {
	TypeImpl
	Element types.Option[Type]
	End     types.Option[Type]
}
