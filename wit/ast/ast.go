/*
package ast models the wit ast
see https://github.com/WebAssembly/component-model/blob/main/design/mvp/WIT.md

Unions are modeled by using union specific interfaces where each member of the union
gets a unique interface. This allows reuse of the member types without polluting the
member type with all the locations where it is used.

example:
--------

myunion   ::= something | other | last
last  	 ::= 'last' id
something ::= 'something' id
other 	 ::= 'other' id
id        ::= \w+

	type MyUnion interface {
		myunion()
	}

	type MyUnionSomething interface {
		MyUnion
		Something() Something
	}

	type MyUnionLast interface {
		MyUnion
		Last() Last
	}

	type MyUnionOther interface {
		MyUnion
		Other() Other
	}

	type Last interface {
		last()
		Id() string
	}

	type Other interface {
		Id() Id
	}

	type Something interface {
		Id() Id
	}

	type Id string
*/
package ast

import "github.com/patrickhuber/go-types"

type Node interface {
	node()
}

type Ast interface {
	Node
	PackageName() types.Option[PackageName]
	Items() []AstItem
}

type PackageName interface {
	Namespace() Id
	Name() Id
	Version() types.Option[Version]
}

type AstItem interface {
	astItem()
}

type AstItemUse interface {
	AstItem
	Use() TopLevelUse
}

type TopLevelUse interface {
	Item() UsePath
	As() types.Option[Id]
}

type UsePath interface {
	usePath()
}

type UsePathId interface {
	UsePath
	Id() Id
}

type UsePathPackage interface {
	UsePath
	Id() PackageName
	Name() Id
}

type Interface interface {
	iface()
	Id() Id
}

type InterfaceFullName interface {
	Interface
	Namespace() Id
	Type() Id
	Version() Version
}

type AstItemInterface interface {
	AstItem
	Interface() Interface
}

type AstItemWorld interface {
	AstItem
	World() World
}

type World interface {
	Name() Id
	Items() []WorldItem
}

type WorldItem interface {
	worldItem()
}

type WorldItemExport interface {
	WorldItem
	Export() Export
}

type Export interface {
	exportItem()
}

type ExportInterface interface {
	Export
	Interface() Interface
}

type ExportExternType interface {
	Export
	Id() Id
	ExternType() ExternType
}

type WorldItemImport interface {
	WorldItem
	Import() Import
}

type Import interface {
	importitem()
}

type ImportInterface interface {
	Import
	Interface() Interface
}

type ImportExternType interface {
	Import
	Id() Id
	ExternType() ExternType
}

type ExternType interface {
	FuncType() FuncType
	InterfaceItems() []InterfaceItems
}

type FuncType interface{}
type InterfaceItems interface{}

type WorldItemUse interface {
	WorldItem
}

type WorldItemTypeDef interface {
	WorldItem
}

type TypeDef interface {
	Name() Id
	Type() Type
}

type WorldItemInclude interface {
	WorldItem
}

type Use interface {
	From() UsePath
	Names() []UseName
}

type UseName interface {
	Name() Id
	As() types.Option[Id]
}

type Type interface {
	ty()
}

type BoolType interface {
	Type
}

type U8Type interface {
	Type
}

type U16Type interface {
	Type
}

type U32Type interface {
	Type
}

type U64Type interface {
	Type
}

type S8Type interface {
	Type
}

type S16Type interface {
	Type
}

type S32Type interface {
	Type
}

type S64Type interface {
	Type
}

type Float32Type interface {
	Type
}

type Float64Type interface {
	Type
}

type CharType interface {
	Type
}

type StringType interface {
	Type
}

type NameType interface {
	Type
	Id() Id
}

type ListType interface {
	Type
	Types() []Type
}

type HandleType interface {
	Type
	Handle() Handle
}

type Handle interface {
	handle()
}

type HandleOwn interface {
	Handle
	Id() Id
}

type HandleBorrow interface {
	Handle
	Id() Id
}

type ResourceType interface {
	Type
	Resource() Resource
}

type Resource interface {
	Functions() []ResourceFunc
}

type ResourceFunc interface {
	resourceFunc()
}

type ResourceFuncMethod interface {
	ResourceFunc
	NamedFunc() NamedFunc
}

type ResourceFuncStatic interface {
	ResourceFunc
	NamedFunc() NamedFunc
}

type ResourceFuncConstructor interface {
	ResourceFunc
	NamedFunc() NamedFunc
}

type NamedFunc interface {
	Id() Id
	Func() Func
}

type Func interface {
	Params() []Param
	Results() ResultList
}

type Param interface {
	Id() Id
	Type() Type
}

type ResultList interface {
	Named() []Param
	Anon() Type
}

type RecordType interface {
	Type
	Record() Record
}

type Record interface{}

type FlagsType interface {
	Type
	Flags() Flags
}

type Flags interface{}

type VariantType interface {
	Type
	Variant() Variant
}

type Variant interface{}

type TupleType interface {
	Type
	Tuple() Tuple
}

type Tuple interface{}

type EnumType interface {
	Type
	Enum() Enum
}

type Enum interface{}

type OptionType interface {
	Type
	Types() Type
}

type ResultType interface {
	Type
	Result() Result
}

type Result interface{}

type FutureType interface {
	Type
	Types() []types.Option[Type]
}

type StreamType interface {
	Type
	Stream() Stream
}

type Stream interface{}

type UnionType interface {
	Type
	Union() Union
}

type Union interface{}

type Id interface{}
type Version interface{}
