package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterIDsMarshalAsStrings(t *testing.T) {
	const big = 516687706984677377

	out, err := json.Marshal(ModuleMappingFilter{NodeID: big, ComposeModuleID: big, ComposeNamespaceID: big, FederationModuleID: big})
	require.NoError(t, err)
	require.Contains(t, string(out), `"nodeID":"516687706984677377","composeModuleID":"516687706984677377","composeNamespaceID":"516687706984677377","federationModuleID":"516687706984677377"`)

	out, err = json.Marshal(NodeSyncFilter{NodeID: big, RelNodeID: big, ModuleID: big})
	require.NoError(t, err)
	require.Contains(t, string(out), `"nodeID":"516687706984677377","relNodeID":"516687706984677377","moduleID":"516687706984677377"`)
}
