package codegen

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestOpenAPIDocsAreCurrent holds server/docs/*.yaml to the rest.yaml files
// they are generated from. UPDATE_DOCS=1 rewrites them, which is also
// what make codegen does on every run.
func TestOpenAPIDocsAreCurrent(t *testing.T) {
	dd, err := procRest(glob("../../*/rest.yaml")...)
	require.NoError(t, err)
	require.NotEmpty(t, dd)

	for _, d := range dd {
		got, err := renderOpenAPI(d)
		require.NoError(t, err)
		p := openAPIPath(d)

		if os.Getenv("UPDATE_DOCS") == "1" {
			require.NoError(t, os.WriteFile(p, got, 0o644))
			continue
		}
		want, err := os.ReadFile(p)
		require.NoErrorf(t, err, "%s is missing; run make codegen", p)
		require.Equalf(t, string(want), string(got), "%s is out of date with %s/rest.yaml; run make codegen and commit the result", path.Base(p), d.App)
	}
}

// TestOpenAPISchemaMapping pins how rest.yaml types read on the wire.
func TestOpenAPISchemaMapping(t *testing.T) {
	require.Equal(t, oaSchema{Type: "string", Format: "uint64", Description: "64-bit integer as a string"}, openAPISchema("uint64"))
	require.Equal(t, "array", openAPISchema("[]string").Type)
	require.Equal(t, "string", openAPISchema("[]uint64").Items.Type)
	require.Equal(t, "date-time", openAPISchema("*time.Time").Format)
	require.Equal(t, "binary", openAPISchema("*multipart.FileHeader").Format)
	require.Equal(t, "object", openAPISchema("types.RecordValueSet").Type)
	require.Equal(t, "boolean", openAPISchema("bool").Type)
}
