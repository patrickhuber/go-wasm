package ast

import "github.com/patrickhuber/go-types"

type PackageNameNode struct {
	namespace Id
	name      Id
	version   types.Option[Version]
}

func (p *PackageNameNode) Namespace() Id {
	return p.namespace
}

func (p *PackageNameNode) Name() Id {
	return p.name
}

func (p *PackageNameNode) Version() types.Option[Version] {
	return p.version
}
