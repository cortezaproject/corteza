package types

import "strconv"

// A little helper to extract identifiers for a given module
func identifiersForModule(m *ComposeModule) NodeIdentifiers {
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
