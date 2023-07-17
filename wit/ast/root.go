package ast

import "github.com/patrickhuber/go-types"

type AstNode struct {
	packageName types.Option[PackageName]
	items       []AstItem
}

func (n *AstNode) PackageName() types.Option[PackageName] {
	return n.packageName
}

func (n *AstNode) Items() []AstItem {
	return n.items
}

func (n *AstNode) node() {}

type AstItemNode struct {
}

func (n *AstItemNode) astItem() {}
func (n *AstItemNode) node()    {}

type AstItemUseNode struct {
	AstItemNode
	topLevelUse TopLevelUse
}

func (n *AstItemUseNode) Use() TopLevelUse {
	return n.topLevelUse
}

type AstItemInterfaceNode struct {
	AstItemNode
	iface Interface
}

func (n *AstItemInterfaceNode) Interface() Interface {
	return n.iface
}

type AstItemWorldNode struct {
	AstItemNode
	worldItem WorldItem
}

func (n *AstItemWorldNode) World() WorldItem {
	return n.worldItem
}
