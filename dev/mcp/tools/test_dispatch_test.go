package tools

import (
	"context"
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

// TestRunnerForNodeSuites pins the one case the suffix has to settle: dev/mcp
// is a Go module that also holds .mjs suites, so a target under it matches the
// Go prefix and would run `go test` on a file go tooling cannot see.
func TestRunnerForNodeSuites(t *testing.T) {
	cases := []struct {
		target string
		want   string
	}{
		{"dev/mcp/shots.test.mjs", "node"},
		{"dev/mcp", "go"},
		{"dev/mcp/tools", "go"},
		{"server/store/adapters/rdbms", "go"},
		{"lib/vue", "vitest"},
	}

	for _, c := range cases {
		t.Run(c.target, func(t *testing.T) {
			if got := runnerFor(t.TempDir(), c.target); got != c.want {
				t.Fatalf("runnerFor(%q) = %q, want %q", c.target, got, c.want)
			}
		})
	}
}

// TestRunNodeReportsFailures runs node's real test runner over a fixture, so
// the TAP parsing is exercised rather than assumed.
func TestRunNodeReportsFailures(t *testing.T) {
	root := t.TempDir()
	suite := "fixture.test.mjs"
	body := `
import { test } from 'node:test'
import assert from 'node:assert'
test('this one holds', () => assert.equal(1, 1))
test('this one does not', () => assert.equal(1, 2))
`
	if err := os.WriteFile(filepath.Join(root, suite), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	report, err := runNode(context.Background(), root, suite, "")
	if err != nil {
		t.Fatalf("runNode: %v", err)
	}
	if report.Passed {
		t.Fatal("report says passed; the fixture has a failing test")
	}
	if len(report.Failures) != 1 || report.Failures[0].Test != "this one does not" {
		t.Fatalf("failures = %+v, want the one failing test named", report.Failures)
	}

	// The green half: filtering to the passing test must report a pass.
	green, err := runNode(context.Background(), root, suite, "this one holds")
	if err != nil {
		t.Fatalf("runNode filtered: %v", err)
	}
	if !green.Passed {
		t.Fatalf("filtered run should pass, got %+v", green.Failures)
	}
}
