package ast

type WorldNode struct {
	name  Id
	items []WorldItem
}

func (n *WorldNode) Name() Id {
	return n.name
}

func (n *WorldNode) Items() []WorldItem {
	return n.items
}

type WorldItemNode struct {
}

func (n *WorldItemNode) worldItem() {}

type WorldItemImportNode struct {
	WorldItemNode
	_import Import
}

func (n *WorldItemImportNode) Import() Import {
	return n._import
}

type WorldItemExportNode struct {
	WorldItemNode
	export Export
}

func (n *WorldItemExportNode) Export() Export {
	return n.export
}

type WorldItemUseNode struct {
	WorldItemNode
	use Use
}

func (n *WorldItemUseNode) Use() Use {
	return n.use
}

type WorldItemTypeDefNode struct {
	WorldItemNode
	typeDef TypeDef
}

func (n *WorldItemTypeDefNode) TypeDef() TypeDef {
	return n.typeDef
}
