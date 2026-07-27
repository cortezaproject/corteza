package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Regression: legacy functions define parameters by Name (ArgumentName empty)
// and the editor sends arguments keyed by Target. Verification must pair them
// on that fallback key instead of collapsing every param onto the first
// empty-ArgumentName argument (which mismatched heterogeneous arg types, e.g.
// composeRecordsLookup namespace/module/record).
func TestParamSet_VerifyArguments_LegacyNameTargetKeying(t *testing.T) {
	set := ParamSet{
		{Name: "namespace", Types: []string{"ID", "Handle", "ComposeNamespace"}, Required: true},
		{Name: "module", Types: []string{"ID", "Handle", "ComposeModule"}, Required: true},
		{Name: "record", Types: []string{"ID", "ComposeRecord"}},
	}

	args := ExprSet{
		{Target: "namespace", Type: "ComposeNamespace", Expr: "namespace"},
		{Target: "module", Type: "ComposeModule", Expr: "module"},
		{Target: "record", Type: "ComposeRecord", Expr: "record"},
	}

	require.NoError(t, set.VerifyArguments(args))
}

// NG automation keys both params and arguments by ArgumentName; that path must
// keep working.
func TestParamSet_VerifyArguments_NGArgumentNameKeying(t *testing.T) {
	set := ParamSet{
		{ArgumentName: "namespace", Types: []string{"ComposeNamespace"}, Required: true},
		{ArgumentName: "module", Types: []string{"ComposeModule"}, Required: true},
	}

	args := ExprSet{
		{ArgumentName: "namespace", Type: "ComposeNamespace"},
		{ArgumentName: "module", Type: "ComposeModule"},
	}

	require.NoError(t, set.VerifyArguments(args))
}

func TestParamSet_VerifyArguments_TypeMismatchStillCaught(t *testing.T) {
	set := ParamSet{
		{Name: "module", Types: []string{"ID", "Handle", "ComposeModule"}, Required: true},
	}

	// Wrong type on the module argument must still fail.
	err := set.VerifyArguments(ExprSet{{Target: "module", Type: "ComposeNamespace"}})
	require.Error(t, err)
}

func TestParamSet_VerifyArguments_RequiredMissing(t *testing.T) {
	set := ParamSet{
		{Name: "module", Types: []string{"ComposeModule"}, Required: true},
	}

	err := set.VerifyArguments(ExprSet{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "module")
}
