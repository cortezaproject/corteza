package types

import (
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	// ComposeModuleNode represents a ComposeModule
	ComposeModuleNode struct {
		Mod *ComposeModule

		// Related namespace
		Ns *ComposeNamespace
	}
)

func (n *ComposeModuleNode) Identifiers() NodeIdentifiers {
	ii := make(NodeIdentifiers, 0)

	if n.Mod.Handle != "" {
		ii = ii.Add(n.Mod.Handle)
	}

	if n.Mod.Name != "" {
		ii = ii.Add(n.Mod.Name)
	}

	if n.Mod.ID > 0 {
		ii = ii.Add(strconv.FormatUint(n.Mod.ID, 10))
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
	if n.Ns.Slug != "" {
		rel.Add(nsr, n.Ns.Slug)
	}
	if n.Ns.Name != "" {
		rel.Add(nsr, n.Ns.Name)
	}
	if n.Ns.ID > 0 {
		rel.Add(nsr, strconv.FormatUint(n.Ns.ID, 10))
	}

	// Related modules via Record module fields
	mdr := types.ModuleRBACResource.String()
	for _, f := range n.Mod.Fields {
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
	n.Ns = mn.Ns
}

func (n *ComposeModuleNode) updateRecFields(m Node) {
	if m.Resource() != types.ModuleRBACResource.String() {
		return
	}
	mn := m.(*ComposeModuleNode)

	// Check what record module field we can link this to
	for _, f := range n.Mod.Fields {
		if f.Kind != "Record" {
			continue
		}

		if mn.Identifiers().HasAny(f.Options.String("module")) {
			f.Options["module"] = strconv.FormatUint(mn.Mod.ID, 10)
		}
	}

	n.Ns = mn.Ns
}
