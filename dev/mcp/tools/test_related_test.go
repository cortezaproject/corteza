package tools

import (
	"reflect"
	"testing"
)

// TestGoRelatedPkgs pins how source files become go test packages: one entry
// per directory, in the order first seen, and nothing from outside the module.
// A file outside the module silently becoming "./" would run the whole module
// and call that a baseline.
func TestGoRelatedPkgs(t *testing.T) {
	got := goRelatedPkgs("server", "server/compose/service/record.go server/compose/service/module.go ./server/pkg/mcpkit/x.go lib/js/src/a.ts")
	want := []string{"./compose/service", "./pkg/mcpkit"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("goRelatedPkgs = %v, want %v", got, want)
	}
	if got := goRelatedPkgs("server", "lib/js/src/a.ts"); len(got) != 0 {
		t.Fatalf("a file outside the module must yield nothing, got %v", got)
	}
}

// TestVitestRelatedArgv pins the related invocation: absolute files, run mode,
// the JSON reporter the parser reads, and --passWithNoTests so an uncovered
// file is a report rather than an exit code.
func TestVitestRelatedArgv(t *testing.T) {
	got := vitestRelatedArgv("/repo", "client/web/unify/src/a.js ./lib/vue/src/b.ts")
	want := []string{"vitest", "related", "/repo/client/web/unify/src/a.js", "/repo/lib/vue/src/b.ts", "--run", "--reporter=json", "--passWithNoTests"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("vitestRelatedArgv = %v, want %v", got, want)
	}
}
