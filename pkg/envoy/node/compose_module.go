package node

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/handle"
	"strconv"
)

type (
	// ComposeModule represents a ComposeModule
	ComposeModule struct {
		Res *types.Module

		// Referenced namespace
		RefNamespace string
	}
)

func (ComposeModule) Resource() string { return "composeModule" }
func (n ComposeModule) Identifiers() envoy.NodeIdentifiers {
	return moduleIdentifiers(n.Res)
}

func (n ComposeModule) Matches(resource string, identifiers ...string) bool {
	if resource != n.Resource() {
		return false
	}

	return n.Identifiers().HasAny(identifiers...)
}

func (n ComposeModule) Relations() envoy.NodeRelationships {
	rel := make(envoy.NodeRelationships)

	// Related namespace
	rel.Add(ComposeNamespace{}.Resource(), identifiers(n.RefNamespace, n.Res.ID)...)

	// Related modules via Record module fields
	for _, f := range n.Res.Fields {
		// @todo should a missing module property raise an error?
		if f.Kind == "Record" && f.Options.String("module") != "" {
			rel.Add(ComposeModule{}.Resource(), f.Options.String("module"))
		}
	}

	return rel
}

func (n *ComposeModule) Update(rr ...envoy.Node) {
	for _, r := range rr {
		switch r := r.(type) {
		case *ComposeNamespace:
			n.updateNamespace(r)

		case *ComposeModule:
			n.updateRecordFields(r)
		}
	}
}

func (n *ComposeModule) updateNamespace(r *ComposeNamespace) {
	if n.RefNamespace == r.Res.Slug {
		n.Res.NamespaceID = r.Res.ID
		n.RefNamespace = ""
	}
}

func (n *ComposeModule) updateRecordFields(r *ComposeModule) {
	for _, f := range n.Res.Fields {
		// here we care only about record field type,
		// ignoring the rest
		if f.Kind != "Record" {
			continue
		}

		// options can store arbitrary values so we can use same opt prop
		// to store module ref (exported) and ID (internal)
		var ref = f.Options.String("module")

		if !handle.IsValid(ref) {
			// empty string or resolved ID, either way,
			// nothing for us to do here..
			continue
		}

		if r.Identifiers().HasAny(ref) {
			f.Options["module"] = strconv.FormatUint(r.Res.ID, 10)
		}
	}
}

// utility fn to extract identifiers for a given module
func moduleIdentifiers(m *types.Module) envoy.NodeIdentifiers {
	return makeIdentifiers(m.Handle, m.ID)
}
