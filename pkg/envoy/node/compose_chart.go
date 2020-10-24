package node

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
)

type (
	// ComposeChart represents a ComposeChart
	ComposeChart struct {
		Res *types.Chart

		// Referenced namespace
		RefNamespace string

		// Referenced modules (via chart reports)
		RefReportModules map[string][]int
	}
)

func (ComposeChart) Resource() string { return "composeChart" }
func (n ComposeChart) Identifiers() envoy.NodeIdentifiers {
	return makeIdentifiers(n.Res.Handle, n.Res.ID)
}

func (n ComposeChart) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n ComposeChart) Relations() envoy.NodeRelationships {
	ref := make(envoy.NodeRelationships)

	// Related namespace
	ref.Add(ComposeNamespace{}.Resource(), identifiers(n.RefNamespace, n.Res.ID)...)

	// Related modules from configure reports
	for modHandle, reports := range n.RefReportModules {
		ref.Add(ComposeModule{}.Resource(), modHandle)
		for _, rIndex := range reports {
			ref.Add(ComposeModule{}.Resource(), identifiers(n.Res.Config.Reports[rIndex].ModuleID)...)
		}
	}

	return ref
}

func (n *ComposeChart) Update(rr ...envoy.Node) {
	for _, r := range rr {
		switch r := r.(type) {
		case *ComposeNamespace:
			n.updateNamespace(r)

		case *ComposeModule:
			n.updateModules(r)
		}
	}
}

func (n *ComposeChart) updateNamespace(r *ComposeNamespace) {
	if n.RefNamespace == r.Res.Slug {
		n.Res.NamespaceID = r.Res.ID
		n.RefNamespace = ""
	}
}

func (n *ComposeChart) updateModules(r *ComposeModule) {
	if reports, has := n.RefReportModules[r.Res.Handle]; has {
		// update module ID on all reports
		for _, rIndex := range reports {
			n.Res.Config.Reports[rIndex].ModuleID = r.Res.ID
		}

		// remove reference
		delete(n.RefReportModules, r.Res.Handle)
	}
}
