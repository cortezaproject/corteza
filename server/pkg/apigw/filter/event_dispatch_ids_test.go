package filter

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapBodyKeepsNumericIDs(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"record":{"id":516687706984677377},"count":3}`))

	vars, err := mapBody(r, map[string][]string{
		"recordID": {"record", "id"},
		"count":    {"count"},
	})
	require.NoError(t, err)
	require.Equal(t, "516687706984677377", vars["recordID"])
	require.Equal(t, float64(3), vars["count"])
}
