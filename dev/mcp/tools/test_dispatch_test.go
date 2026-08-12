package tools

import (
	"os"
	"path/filepath"
	"testing"
)

// TestRunnerDispatch pins which runner a target gets.
//
// lib/js runs mocha and the rest of the JS workspaces run vitest. Sending lib/js
// to vitest is not a loud failure: vitest matches no files, exits clean, and the
// tool reports "no tests ran: nothing matched lib/js — check the path", which
// reads as a mistyped target. The path was right; the runner was wrong, and the
// caller pays for that by hunting for a spec glob that does not exist.
func TestRunnerDispatch(t *testing.T) {
	root := t.TempDir()

	// Only lib/js carries a mocha config, which is what decides this.
	mustWrite(t, filepath.Join(root, "lib", "js", ".mocharc.js"), "module.exports = { bail: true }\n")
	mustWrite(t, filepath.Join(root, "lib", "vue", "package.json"), "{}\n")

	cases := []struct {
		target string
		want   string
	}{
		{"lib/js", "mocha"},
		{"lib/js/src/compose/types/chart/empty-dimension.test.ts", "mocha"},
		{"lib/vue", "vitest"},
		{"client/web/unify", "vitest"},
		{"server/compose/service", "go"},
		{"server/...", "go"},
		{"dev/mcp/tools", "go"},
	}

	for _, c := range cases {
		t.Run(c.target, func(t *testing.T) {
			if got := runnerFor(root, c.target); got != c.want {
				t.Fatalf("runnerFor(%q) = %q, want %q", c.target, got, c.want)
			}
		})
	}
}

// TestMochaBailsIsReportedWhenConfigured covers the caveat attached to a failing
// mocha run: bail stops at the first failure, so the tests after it never ran
// and are absent rather than passing. Reporting one failure without saying so
// reads as "the rest were fine".
func TestMochaBailsIsReportedWhenConfigured(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"js config with bail", "module.exports = {\n  bail: true,\n}\n", true},
		{"js config without bail", "module.exports = {\n  recursive: true,\n}\n", false},
		{"json config with bail", `{ "bail": true }`, true},
		{"json config with bail off", `{ "bail": false }`, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".mocharc.js")
			mustWrite(t, path, c.content)

			if got := mochaBails(path); got != c.want {
				t.Fatalf("mochaBails() = %v, want %v", got, c.want)
			}
		})
	}

	if mochaBails("") {
		t.Fatal("mochaBails(\"\") = true, want false — no config means no caveat")
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
