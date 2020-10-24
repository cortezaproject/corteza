package node

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy"
)

type (
	ComposeRecord struct {
		Res *types.Record

		// resolved module, we'll need access to this when we'll be building
		// relationships
		module *types.Module

		// module reference
		RefModule string

		// all base user references
		RefUsers []string
	}
)

func (n *ComposeRecord) Resource() string { return "composeRecord" }
func (n *ComposeRecord) Identifiers() envoy.NodeIdentifiers {
	return nil
}

func (n *ComposeRecord) Matches(resource string, identifiers ...string) bool {
	return resource == n.Resource() && n.Identifiers().HasAny(identifiers...)
}

func (n *ComposeRecord) Relations() envoy.NodeRelationships {
	if n.module == nil {
		return envoy.NodeRelationships{}
	}

	// This omits the namespace rel. as it's transitively implied via the modules

	rel := make(envoy.NodeRelationships)
	//mdr := types.ModuleRBACResource.String()
	//rrr := "compose:record:"
	//
	//if n.Mod == nil {
	//	return rel
	//}

	//// Original module
	//mIdentifiers := moduleIdentifiers(n.Mod)
	//rel.Add(mdr, mIdentifiers...)
	//
	//// Field relationships
	for _, f := range n.module.Fields {
		switch f.Kind {
		case "Record":
			modID := f.Options.String("module")
			// For the module
			rel.Add(ComposeModule{}.Resource(), modID)

			// For the records.
			// Since this record depends on another module, it also depends on those records.
			rel.Add(n.Resource(), modID)
		}
	}

	return rel
}

func (n *ComposeRecord) Update(mm ...envoy.Node) {
	for _, m := range mm {
		switch c := m.(type) {
		case *ComposeModule:
			// Direct module dependency
			if n.Matches(n.Resource(), moduleIdentifiers(c.Res)...) {
				n.module = c.Res
			}
		}
	}
}
