package reporting

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefinitionID(t *testing.T) {
	var def map[string]any
	require.NoError(t, json.Unmarshal([]byte(`{"moduleID":516687706984677377,"namespaceID":"516687706984153089","connectionID":0}`), &def))

	_, err := definitionID("moduleID", def["moduleID"])
	require.ErrorContains(t, err, "moduleID is a number too large")

	id, err := definitionID("namespaceID", def["namespaceID"])
	require.NoError(t, err)
	require.Equal(t, uint64(516687706984153089), id)

	id, err = definitionID("connectionID", def["connectionID"])
	require.NoError(t, err)
	require.Zero(t, id)
}
