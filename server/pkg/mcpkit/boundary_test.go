package mcpkit

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNoDomainImports is the constraint the whole package exists to hold.
//
// mcpkit is shared between the configurator MCP, which is product code, and the
// developer MCP under dev/, which is a separate Go module that must not drag
// Human's domain — or its vendor tree — along with it. A single import of a
// compose or system package would make that module transitively depend on the
// server's entire object graph, and the split would have bought nothing.
//
// Domain packages may depend on mcpkit. mcpkit may depend on nothing of theirs.
// pkg/ is allowed: it is the shared-plumbing tier both sides already use.
func TestNoDomainImports(t *testing.T) {
	const (
		module  = "github.com/crusttech/human/server/"
		allowed = "pkg/"
	)

	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no source files found — has the package moved?")
	}

	fset := token.NewFileSet()
	for _, name := range files {
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}

		f, err := parser.ParseFile(fset, name, src, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}

		for _, imp := range f.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(path, module) {
				continue
			}

			rel := strings.TrimPrefix(path, module)
			if strings.HasPrefix(rel, allowed) {
				continue
			}

			t.Errorf(
				"%s imports %s: mcpkit is shared with the developer MCP and may not depend on Human's domain. "+
					"If this package needs something from %s, invert it — expose a neutral seam here and let the "+
					"domain side project onto it, the way Registry.Select does for the agentic runtime",
				name, path, rel,
			)
		}
	}
}
