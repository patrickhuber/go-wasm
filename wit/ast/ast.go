package ast

import (
	"github.com/patrickhuber/go-types"
	abi "github.com/patrickhuber/go-wasm/abi/types"
)

type Ast struct {
	PackageName types.Option[PackageName]
	Items       []AstItem
}

type PackageName struct {
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

type InterfaceItem struct {
	Func *NamedFunc
}

type NamedFunc struct {
	Name string
	Func *Func
}

type Func struct {
	Params  []Parameter
	Results []Result
}

type Parameter struct {
	Id   string
	Type abi.Type
}

type Result struct{}

type World struct{}
type TopLevelUse struct{}
