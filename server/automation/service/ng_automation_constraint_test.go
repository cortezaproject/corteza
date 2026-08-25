package service

import (
	"testing"

	"github.com/crusttech/human/server/automation/types"
	"github.com/stretchr/testify/require"
)

func constraint(name, typ, val string) types.NgTriggerConstraint {
	return types.NgTriggerConstraint{
		Name:   name,
		Op:     "eq",
		Values: []types.NgTriggerConstraintValue{{Type: typ, Value: val}},
	}
}

// The construct library advertises the Compose* resource types
// (legacy_triggers.go recordConstraints), so a constraint written from that
// contract has to build rather than be rejected.
func TestPrepConstraintBitsResourceTypes(t *testing.T) {
	svc := &ngAutomation{}

	for _, tc := range []struct {
		typ  string
		want string
	}{
		{"String", "module.name"},
		{"Handle", "module.handle"},
		{"ID", "module.id"},
		{"ComposeNamespace", "module.id"},
		{"ComposeModule", "module.id"},
		{"ComposeRecord", "module.id"},
	} {
		t.Run(tc.typ, func(t *testing.T) {
			name, op, values, err := svc.prepConstraintBits(constraint("module", tc.typ, "tickets"))
			require.NoError(t, err)
			require.Equal(t, tc.want, name)
			require.Equal(t, "eq", op)
			require.Equal(t, []string{"tickets"}, values)
		})
	}
}

// An unusable constraint must be an error the caller can refuse on. Returning
// no error here is what let registerTrigger drop the constraint and register
// the trigger more broadly than it was authored.
func TestPrepConstraintBitsUnknownTypeErrors(t *testing.T) {
	svc := &ngAutomation{}

	_, _, _, err := svc.prepConstraintBits(constraint("module", "NotAType", "tickets"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unknown constraint type")
}
