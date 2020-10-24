package node

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
)

type (
	ComposeNamespace struct {
		Res *types.Namespace
	}
)

func (ComposeNamespace) Resource() string { return "composeNamespace" }
func (n ComposeNamespace) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.Slug, n.Res.ID)
}

func (n ComposeNamespace) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n ComposeNamespace) Relations() envoy.NodeRelationships {
	return nil
}
