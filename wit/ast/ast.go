package ast

import (
	"github.com/patrickhuber/go-types"
)

type Ast struct {
	PackageDeclaration types.Option[PackageDeclaration]
	Items              []AstItem
}

type PackageDeclaration struct {
	Namespace string
	Name      string
	Version   types.Option[Version]
}

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
	Pre   string
	Build string
}

type AstItem struct {
	Interface *Interface
	World     *World
	Use       *TopLevelUse
}

type Interface struct {
	Name  string
	Items []InterfaceItem
}

type InterfaceItem interface {
	interfaceItem()
}

type FuncItem struct {
	InterfaceItem
	ID       string
	FuncType *FuncType
}

type FuncType struct {
	Params  []Parameter
	Results *ResultList
}

type ResultList struct {
	Named     []Parameter
	Anonymous Type
}

type Parameter struct {
	Id   string
	Type Type
}

type World struct {
	Id    string
	Items []WorldItem
}

type WorldItem interface {
	worldItem()
}

type Export struct {
	WorldItem
	ExternType ExternType
}

type Import struct {
	WorldItem
	ExternType ExternType
}

type ExternType interface {
	externType()
}

type ExternTypeFunc struct {
	ExternType
	ID   string
	Func *FuncType
}

type ExternTypeInterface struct {
	ExternType
	ID             string
	InterfaceItems []InterfaceItem
}

type ExternTypeUsePath struct {
	ExternType
	UsePath *UsePath
}

type Use struct {
	WorldItem
	InterfaceItem
	From  *UsePath
	Names []UseName
}

type UsePath struct {
	Id      string
	Package struct {
		Id   *PackageDeclaration
		Name string
	}
}

type UseName struct {
	Name string
	As   types.Option[string]
}

type TypeDef interface {
	WorldItem
	InterfaceItem
	typedef()
}

type Include struct {
	WorldItem
	From  *UsePath
	Names []IncludeName
}

type IncludeName struct {
	Name string
	As   string
}

type TopLevelUse struct {
	Item *UsePath
	As   types.Option[string]
}

type Type interface {
	ty()
}

type U8 struct{ Type }
type U16 struct{ Type }
type U32 struct{ Type }
type U64 struct{ Type }
type S8 struct{ Type }
type S16 struct{ Type }
type S32 struct{ Type }
type S64 struct{ Type }
type Float32 struct{ Type }
type Float64 struct{ Type }
type Char struct{ Type }
type Bool struct{ Type }
type String struct{ Type }

type Tuple struct {
	Type
	Types []Type
}

type List struct {
	Type
	ItemType Type
}

type Option struct {
	Type
	ItemType Type
}
type Result struct {
	Type
	Ok    types.Option[Type]
	Error types.Option[Type]
}

type Handle interface {
	handle()
}

type Own struct {
	Handle
	Type
	Id string
}

type Borrow struct {
	Handle
	Type
	Id string
}

type Id struct {
	Type
	Value string
}

type Stream struct {
	Type
	Element types.Option[Type]
	End     types.Option[Type]
}

type Resource struct {
	TypeDef
	ID      string
	Methods []ResourceMethod
}

type ResourceMethod interface {
	resourceMethod()
}

type Static struct {
	ResourceMethod
	ID       string
	FuncType *FuncType
}

type Constructor struct {
	ResourceMethod
	ParameterList []Parameter
}

type Method struct {
	ResourceMethod
	ID   string
	Func *FuncItem
}

type Record struct {
	TypeDef
	ID     string
	Fields []Field
}

type Field struct {
	Name string
	Type Type
}

type Flags struct {
	TypeDef
	ID    string
	Flags []Flag
}

type Flag struct {
	Id string
}

type Variant struct {
	TypeDef
	ID    string
	Cases []Case
}

type Case struct {
	Name string
	Type types.Option[Type]
}

type Enum struct {
	TypeDef
	ID    string
	Cases []EnumCase
}

type EnumCase struct {
	Name string
}

type TypeItem struct {
	TypeDef
	ID   string
	Type Type
}

type Future struct {
	Type
	ItemType types.Option[Type]
}
