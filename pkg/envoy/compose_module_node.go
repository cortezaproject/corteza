package envoy

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposeModuleNode represents a ComposeModule
	ComposeModuleNode struct {
		Module *types.Module

		// Related namespace
		RefNamespaceSlug string
		RefNamespaceID   uint64
	}
)

func (n *ComposeModuleNode) Identifiers() NodeIdentifiers {
	ii := make(NodeIdentifiers, 0)

	if n.Module.Handle != "" {
		ii = ii.Add(n.Module.Handle)
	}

	if n.Module.Name != "" {
		ii = ii.Add(n.Module.Name)
	}

	if n.Module.ID > 0 {
		ii = ii.Add(strconv.FormatUint(n.Module.ID, 10))
	}

	return ii
}

func (n *ComposeModuleNode) Matches(resource string, identifiers ...string) bool {
	if resource != n.Resource() {
		return false
	}

	return n.Identifiers().HasAny(identifiers...)
}

func (n *ComposeModuleNode) Resource() string {
	return types.ModuleRBACResource.String()
}

func (n *ComposeModuleNode) Relations() NodeRelationships {
	rel := make(NodeRelationships)

	// Related namespace
	nsr := types.NamespaceRBACResource.String()
	if n.RefNamespaceSlug != "" {
		rel.Add(nsr, n.RefNamespaceSlug)
	}

	if n.RefNamespaceID > 0 {
		rel.Add(nsr, strconv.FormatUint(n.RefNamespaceID, 10))
	}

	// Related modules via Record module fields
	mdr := types.ModuleRBACResource.String()
	for _, f := range n.Module.Fields {
		// @todo should a missing module property raise an error?
		if f.Kind == "Record" && f.Options.String("module") != "" {
			rel.Add(mdr, f.Options.String("module"))
		}
	}

	return rel
}

func (n *ComposeModuleNode) Update(mm ...Node) {
	for _, m := range mm {
		n.updateNamespace(m)
		n.updateRecFields(m)
	}
}

func (n *ComposeModuleNode) updateNamespace(m Node) {
	if m.Resource() != types.NamespaceRBACResource.String() {
		return
	}

	mn := m.(*ComposeNamespaceNode)
	n.Module.NamespaceID = mn.Ns.ID
}

func (n *ComposeModuleNode) updateRecFields(m Node) {
	if m.Resource() != types.ModuleRBACResource.String() {
		return
	}

	mn := m.(*ComposeModuleNode)

	// Check what record module field we can link this to
	for _, f := range n.Module.Fields {
		if f.Kind != "Record" {
			continue
		}

		if mn.Identifiers().HasAny(f.Options.String("module")) {
			f.Options["module"] = strconv.FormatUint(mn.Module.ID, 10)
		}
	}
}
