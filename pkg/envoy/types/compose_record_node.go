package types

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	RecordIterator func(func(*ComposeRecord) error) error

	ComposeRecordNode struct {
		Walk RecordIterator

		// Metafields for relationship management
		Mod *ComposeModule
	}
)

func (n *ComposeRecordNode) Identifiers() NodeIdentifiers {
	return identifiersForModule(n.Mod)
}

func (n *ComposeRecordNode) Resource() string {
	return "compose:record:"
}

func (n *ComposeRecordNode) Matches(resource string, identifiers ...string) bool {
	if resource != n.Resource() {
		return false
	}

	return n.Identifiers().HasAny(identifiers...)
}

func (n *ComposeRecordNode) Relations() NodeRelationships {
	// This omits the namespace rel. as it's transitively implied via the modules

	rel := make(NodeRelationships)
	mdr := types.ModuleRBACResource.String()
	rrr := "compose:record:"

	if n.Mod == nil {
		return rel
	}

	// Original module
	mIdentifiers := identifiersForModule(n.Mod)
	rel.Add(mdr, mIdentifiers...)

	// Field relationships
	for _, f := range n.Mod.Fields {
		if f.Kind == "Record" {
			modID := f.Options.String("module")
			// For the module
			rel.Add(mdr, modID)

			// For the records.
			// Since this record depends on another module, it also depends on those records.
			rel.Add(rrr, modID)
		}
	}

	return rel
}

func (n *ComposeRecordNode) Update(mm ...Node) {
	for _, m := range mm {
		switch m.Resource() {
		case types.ModuleRBACResource.String():
			mn, _ := m.(*ComposeModuleNode)
			// Direct module dependency
			if n.Matches(n.Resource(), identifiersForModule(mn.Mod)...) {
				n.Mod = mn.Mod
			}
		}
	}
}
