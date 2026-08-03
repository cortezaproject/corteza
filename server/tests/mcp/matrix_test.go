package mcp_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// update rewrites the committed matrix instead of asserting against it:
//
//	cd server && go test ./tests/mcp/ -run TestToolsMatrix -update
var update = flag.Bool("update", false, "rewrite TOOLS.md from the current registry")

// matrixPath is where the generated coverage matrix lives — next to the
// registry it describes, so an agent looking for the tool surface finds it at a
// predictable path rather than having to know a build step.
func matrixPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "server", "system", "agentic", "mcp", "TOOLS.md")
}

// TestToolsMatrix keeps TOOLS.md in step with the registry.
//
// A golden file rather than a build step: generation is a side effect of the
// test, so a stale matrix fails CI instead of quietly describing a surface that
// no longer exists.
func TestToolsMatrix(t *testing.T) {
	want := renderMatrix(buildRegistry(t))
	path := matrixPath(t)

	if *update {
		require.NoError(t, os.WriteFile(path, []byte(want), 0o644))
		t.Logf("wrote %s", path)
		return
	}

	got, err := os.ReadFile(path)
	require.NoErrorf(t, err, "TOOLS.md is missing; regenerate with -update")

	require.Equalf(t, want, string(got),
		"TOOLS.md is stale. Regenerate:\n\n"+
			"    cd server && go test ./tests/mcp/ -run TestToolsMatrix -update\n")
}
