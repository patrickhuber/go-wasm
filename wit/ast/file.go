package ast

type FileNode struct {
	packageDeclaration PackageDeclaration
	topLevelItems      []TopLevelItem
}

func (n *FileNode) PackageDeclaration() PackageDeclaration {
	return n.packageDeclaration
}

func (n *FileNode) TopLevelItems() []TopLevelItem {
	return n.topLevelItems
}
