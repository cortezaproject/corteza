package expr

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVarsUnmarshalKeepsNumericIDs(t *testing.T) {
	vars := &Vars{}
	require.NoError(t, json.Unmarshal([]byte(`{
		"id":  {"@type": "ID", "@value": 516687706984677377},
		"int": {"@type": "Integer", "@value": 516687706984677377},
		"n":   {"@type": "Integer", "@value": 42}
	}`), vars))

	id, err := CastToID(vars.GetValue()["id"].Get())
	require.NoError(t, err)
	require.Equal(t, uint64(516687706984677377), id)

	i, err := CastToInteger(vars.GetValue()["int"].Get())
	require.NoError(t, err)
	require.Equal(t, int64(516687706984677377), i)

	require.Equal(t, float64(42), vars.GetValue()["n"].Get())
}
