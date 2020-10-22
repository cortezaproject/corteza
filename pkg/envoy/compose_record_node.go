package envoy

import (
	"github.com/cortezaproject/corteza-server/compose/types"
	"strconv"
)

type (
	RecordIterator func(func(record *types.Record) error) error

	ComposeRecordNode struct {
		Walk RecordIterator

		RefModuleHandle string
	}
)

func (n *ComposeRecordNode) Identifiers() NodeIdentifiers {
	//return identifiersForModule(n.Mod)
	return nil
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
	//mdr := types.ModuleRBACResource.String()
	//rrr := "compose:record:"
	//
	//if n.Mod == nil {
	//	return rel
	//}

	//// Original module
	//mIdentifiers := identifiersForModule(n.Mod)
	//rel.Add(mdr, mIdentifiers...)
	//
	//// Field relationships
	//for _, f := range n.Mod.Fields {
	//	if f.Kind == "Record" {
	//		modID := f.Options.String("module")
	//		// For the module
	//		rel.Add(mdr, modID)
	//
	//		// For the records.
	//		// Since this record depends on another module, it also depends on those records.
	//		rel.Add(rrr, modID)
	//	}
	//}

	return rel
}

func (n *ComposeRecordNode) Update(mm ...Node) {
	for _, m := range mm {
		switch m.Resource() {
		case types.ModuleRBACResource.String():
			mn, _ := m.(*ComposeModuleNode)
			// Direct module dependency
			if n.Matches(n.Resource(), identifiersForModule(mn.Module)...) {
				//n.Mod = mn.Module
			}
		}
	}
}

// A little helper to extract identifiers for a given module
func identifiersForModule(m *types.Module) NodeIdentifiers {
	ii := make(NodeIdentifiers, 0)

	if m == nil {
		return ii
	}

	if m.Handle != "" {
		ii = ii.Add(m.Handle)
	}
	if m.Name != "" {
		ii = ii.Add(m.Name)
	}
	if m.ID > 0 {
		ii = ii.Add(strconv.FormatUint(m.ID, 10))
	}

	return ii
}
