package mcp_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGlossaryNamesEveryToolFamily holds docs/get-started/glossary.md to the
// registry: every tool belongs to a family whose prefix the glossary names in
// backticks, so a term read on screen can be found in the API. A tool family
// added without a glossary row fails here.
func TestGlossaryNamesEveryToolFamily(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join(repoRoot(t), "docs", "get-started", "glossary.md"))
	require.NoError(t, err)

	prefixes := map[string]bool{}
	for _, m := range regexp.MustCompile("`([a-z][a-z0-9_]*)`").FindAllStringSubmatch(string(raw), -1) {
		prefixes[m[1]] = true
	}

	for _, tool := range buildRegistry(t).Tools() {
		named := false
		for p := range prefixes {
			if tool.Name == p || strings.HasPrefix(tool.Name, p+"_") {
				named = true
				break
			}
		}
		require.Truef(t, named, "tool %q belongs to no family the glossary names; add its prefix in backticks to docs/get-started/glossary.md", tool.Name)
	}
}
