package id

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUint64JSON(t *testing.T) {
	out, err := json.Marshal(Uint64(516687706984677377))
	require.NoError(t, err)
	require.Equal(t, `"516687706984677377"`, string(out))

	for in, want := range map[string]Uint64{
		`"516687706984677377"`: 516687706984677377,
		`516687706984677377`:   516687706984677377,
		`""`:                   0,
		`null`:                 0,
	} {
		var v Uint64
		require.NoError(t, json.Unmarshal([]byte(in), &v), in)
		require.Equal(t, want, v, in)
	}

	var v Uint64
	require.Error(t, json.Unmarshal([]byte(`"abc"`), &v))
}

func TestUint64sJSON(t *testing.T) {
	out, err := json.Marshal(Uint64s{516687706984677377, 1})
	require.NoError(t, err)
	require.Equal(t, `["516687706984677377","1"]`, string(out))

	out, err = json.Marshal(Uint64s(nil))
	require.NoError(t, err)
	require.Equal(t, `null`, string(out))

	var vv Uint64s
	require.NoError(t, json.Unmarshal([]byte(`["516687706984677377",1]`), &vv))
	require.Equal(t, Uint64s{516687706984677377, 1}, vv)
}
