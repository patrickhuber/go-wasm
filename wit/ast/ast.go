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

type File interface {
	PackageDeclaration() PackageDeclaration
	TopLevelItems() []TopLevelItem
}

type PackageDeclaration interface {
	Namespace() Id
	Id() Id
	Version() Version
}

type TopLevelItem interface {
	toplevelitem()
}

type TopLevelUse interface {
	TopLevelItem
	Use() Interface
	As() Id
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

type TopLevelIterface interface {
	TopLevelItem
	Interface() Interface
}

type TopLeveWorld interface {
	TopLevelItem
	WorldItem() WorldItem
}

type WorldItem interface {
	World() Id
	Items() []WorldItems
}

type WorldItems interface {
	worlditems()
}

type WorldItemExportItem interface {
	WorldItems
	ExportItem() ExportItem
}

type ExportItem interface {
	exportitem()
}

type ExportItemInterface interface {
	ExportItem
	Interface() Interface
}

type ExportItemExternType interface {
	ExportItem
	Id() Id
	ExternType() ExternType
}

type WorldItemImportItem interface {
	WorldItems
	ImportItem() ImportItem
}

type ImportItem interface {
	importitem()
}

type ImportItemInterface interface {
	ImportItem
	Interface() Interface
}

type ImportItemExternType interface {
	ExportItem
	Id() Id
	ExternType() ExternType
}

type ExternType interface {
	FuncType() FuncType
	InterfaceItems() []InterfaceItems
}

type FuncType interface{}
type InterfaceItems interface{}

type WorldItemUseItem interface {
	WorldItems
}

type WorldItemTypeDefItem interface {
	WorldItems
}

type WorldItemIncludeItem interface {
	WorldItems
}

type Id interface{}
type Version interface{}
