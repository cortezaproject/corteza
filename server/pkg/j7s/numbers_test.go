package j7s

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsafeNumber(t *testing.T) {
	var v any
	require.NoError(t, json.Unmarshal([]byte(`{"a":[1,{"b":516687706984677377}]}`), &v))

	path, found := UnsafeNumber(v)
	require.True(t, found)
	require.Equal(t, "a[1].b", path)

	_, found = UnsafeNumber(map[string]any{"a": float64(MaxExactNumber - 1), "b": "516687706984677377"})
	require.False(t, found)
}

func TestExactNumbers(t *testing.T) {
	var v any
	dec := json.NewDecoder(bytes.NewReader([]byte(`{"id":516687706984677377,"n":42,"f":1.5,"neg":-516687706984677377,"list":[516687706984677377]}`)))
	dec.UseNumber()
	require.NoError(t, dec.Decode(&v))

	require.Equal(t, map[string]any{
		"id":   "516687706984677377",
		"n":    float64(42),
		"f":    1.5,
		"neg":  "-516687706984677377",
		"list": []any{"516687706984677377"},
	}, ExactNumbers(v))
}
