package corredor

import (
	"fmt"
	"strings"
)

type (
	// ScriptSecurity sets script security and run-as flag
	//
	// Determinate user (if >0) that script will run under
	// and allow/deny combination for roles that can see & execute
	// explicit (onManual) scripts
	//
	// Please note that implicit scripts cannot have RBAC protection!
	// These scripts are always executed as configured no mather who runs them.
	ScriptSecurity struct {
		*Security
		runAs uint64
	}

	Script struct {
		Name        string          `json:"name"`
		Label       string          `json:"label"`
		Description string          `json:"description"`
		Errors      []string        `json:"errors,omitempty"`
		Triggers    []*Trigger      `json:"triggers"`
		Iterator    *Iterator       `json:"iterator"`
		Security    *ScriptSecurity `json:"security"`

		// If bundle or type is set, consider
		// this a frontend script
		Bundle string `json:"bundle,omitempty"`
		Type   string `json:"type,omitempty"`
	}
)

// FindByName returns script from the set if it exists
func (set ScriptSet) FindByName(name string) *Script {
	for i := range set {
		if set[i].Name == name {
			return set[i]
		}
	}

	return nil
}

// String fn to make logging a bit friendlier
func (ss *ScriptSecurity) String() (o string) {
	if ss == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"runAs: %s, allow: %s, deny: %s",
		ss.RunAs,
		strings.Join(ss.Allow, ","),
		strings.Join(ss.Deny, ","),
	)
}
