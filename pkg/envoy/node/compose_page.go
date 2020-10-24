package node

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
)

type (
	// ComposePage represents a ComposePage
	ComposePage struct {
		Res *types.Page

		// Referenced namespace
		RefNamespace string

		// Referenced page module
		RefModule string
	}
)

func (ComposePage) Resource() string { return "composePage" }
func (n ComposePage) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.Handle, n.Res.ID)
}

func (n ComposePage) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n *ComposePage) Relations() envoy.NodeRelationships {
	rel := make(envoy.NodeRelationships)

	// Related namespace
	rel.Add(ComposeNamespace{}.Resource(), n.RefNamespace)
	rel.Add(ComposeModule{}.Resource(), n.RefModule)

	// Related modules from page blocks
	for _, block := range n.Res.Blocks {
		// @todo detect block type and collect all references from block options
		_ = block
	}

	return rel
}

func (n *ComposePage) Update(rr ...envoy.Node) {
	for _, r := range rr {
		switch r := r.(type) {
		case *ComposeNamespace:
			n.updateNamespace(r)

		case *ComposeModule:
			n.updateModules(r)

			// @todo any other types of references we need to resolve?
		}
	}
}

func (n *ComposePage) updateNamespace(r *ComposeNamespace) {
	if n.RefNamespace == r.Res.Slug {
		n.Res.NamespaceID = r.Res.ID
		n.RefNamespace = ""
	}
}

func (n *ComposePage) updateModules(r *ComposeModule) {
	if r.Identifiers().HasAny(identifiers(n.RefModule, n.Res.ModuleID)...) {
		n.Res.ModuleID = r.Res.ID
	}
}
