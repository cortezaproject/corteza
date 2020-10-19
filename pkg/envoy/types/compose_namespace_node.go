package types

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeNamespaceNode struct {
		Ns *ComposeNamespace
	}
)

func (n *ComposeNamespaceNode) Identifiers() NodeIdentifiers {
	ii := make(NodeIdentifiers, 0)

	if n.Ns.Slug != "" {
		ii = ii.Add(n.Ns.Slug)
	}
	if n.Ns.Name != "" {
		ii = ii.Add(n.Ns.Name)
	}
	if n.Ns.ID > 0 {
		ii = ii.Add(strconv.FormatUint(n.Ns.ID, 10))
	}

	return ii
}

func (n *ComposeNamespaceNode) Matches(resource string, identifiers ...string) bool {
	if resource != n.Resource() {
		return false
	}

	return n.Identifiers().HasAny(identifiers...)
}

func (n *ComposeNamespaceNode) Resource() string {
	return types.NamespaceRBACResource.String()
}

func (n *ComposeNamespaceNode) Relations() NodeRelationships {
	return nil
}
