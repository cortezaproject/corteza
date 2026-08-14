package tools

import (
	"strings"
	"testing"
)

// TestCheckMessage pins the rules that exist because they were forgotten by
// hand. Each case is a message that reached, or nearly reached, this repo's
// history before the tool existed.
func TestCheckMessage(t *testing.T) {
	cases := []struct {
		name    string
		subject string
		body    string
		want    string // substring of the expected refusal; empty means accept
	}{
		{name: "ordinary subject", subject: "Fix the argument binding path"},
		{name: "subject with a short body", subject: "Resolve notification recipient IDs",
			body: "The ID path trusted its input, so a nonexistent user\nwrote a notification addressed to nobody."},

		{name: "past tense", subject: "Fixed the argument binding path", want: "not imperative"},
		{name: "third person", subject: "Fixes the argument binding path", want: "not imperative"},
		{name: "too short", subject: "Fix bug", want: "too short"},
		{name: "trailing period", subject: "Fix the argument binding path.", want: "period"},
		{name: "lowercase", subject: "fix the argument binding path", want: "capital"},
		{
			name:    "too long",
			subject: "Fix the argument binding path so that every single caller of the converter behaves",
			want:    "over the 72",
		},

		{
			name:    "co-author trailer",
			subject: "Fix the argument binding path",
			body:    "Co-Authored-By: Claude <noreply@anthropic.com>",
			want:    "AI attribution",
		},
		{
			name:    "generated-with line",
			subject: "Fix the argument binding path",
			body:    "🤖 Generated with [Claude Code](https://claude.com/claude-code)",
			want:    "AI attribution",
		},
		{
			name:    "essay body",
			subject: "Fix the argument binding path",
			body:    "one\ntwo\nthree\nfour\nfive\nsix",
			want:    "longer than three",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkMessage(c.subject, c.body)

			if c.want == "" {
				if len(got) != 0 {
					t.Fatalf("expected %q to be accepted, refused with %v", c.subject, got)
				}
				return
			}

			if len(got) == 0 {
				t.Fatalf("expected %q to be refused for %q, accepted", c.subject, c.want)
			}

			if !strings.Contains(strings.Join(got, " "), c.want) {
				t.Errorf("refusal does not mention %q: %v", c.want, got)
			}
		})
	}
}

// TestCheckAtomic covers the pairings that are fine and the ones that are not.
//
// It warns rather than refuses, so what matters is that it stays quiet on the
// combinations people legitimately commit together — otherwise it becomes noise
// to ignore, which is worse than not existing.
func TestCheckAtomic(t *testing.T) {
	quiet := [][]string{
		{"server/automation/service/workflow.go"},
		{"server/automation/service/workflow.go", "server/automation/service/workflow_test.go"},
		{"server/system/types/user.go", "server/system/types/user.gen.go"},
		{"docs/one.md", "docs/two.md"},
		// An intent doc states what the code change beside it made true, so the
		// two belong in one commit — and the lock the sync writes travels with
		// the doc.
		{"server/pkg/mcpkit/registry.go", "server/pkg/mcpkit/mcpkit.intent.md"},
		{"server/pkg/mcpkit/registry.go", "server/pkg/mcpkit/mcpkit.intent.md", ".intent/intent.lock.json"},
		{"server/pkg/mcpkit/mcpkit.intent.md", ".intent/intent.lock.json"},
	}

	for _, files := range quiet {
		if got := checkAtomic(files); len(got) != 0 {
			t.Errorf("expected no warning for %v, got %v", files, got)
		}
	}

	loud := [][]string{
		{"server/automation/service/workflow.go", "README.md"},
		{"server/automation/service/workflow.go", "locale/en/human-webapp/project.yaml"},
		// An intent doc riding along does not make an unrelated prose doc fine.
		{"server/pkg/mcpkit/registry.go", "server/pkg/mcpkit/mcpkit.intent.md", "README.md"},
	}

	for _, files := range loud {
		if got := checkAtomic(files); len(got) == 0 {
			t.Errorf("expected a warning for %v, got none", files)
		}
	}
}

// TestIsIntentPath pins what counts as the intent half of a commit. The lock is
// the easy one to miss: sync writes it, and a commit that leaves it behind
// reports drift that was already reconciled.
func TestIsIntentPath(t *testing.T) {
	yes := []string{
		"server/pkg/mcpkit/mcpkit.intent.md",
		"client/web/unify/src/sections/compose/views/Admin/Modules/Edit.intent.md",
		".intent/intent.lock.json",
		".intent/SPEC.md",
	}

	for _, f := range yes {
		if !isIntentPath(f) {
			t.Errorf("expected %q to count as an intent path", f)
		}
	}

	no := []string{
		"README.md",
		"docs/intent.md",
		"server/automation/service/workflow.go",
		"client/web/unify/src/sections/compose/views/Admin/Modules/Edit.vue",
	}

	for _, f := range no {
		if isIntentPath(f) {
			t.Errorf("expected %q not to count as an intent path", f)
		}
	}
}

// TestCheckIntentPairedStaysQuiet covers the one branch that must not shell out:
// a commit already carrying an intent change needs no drift check, and running
// one would put a node subprocess in the path of every commit that did the right
// thing. A nil root would fail loudly if the short-circuit ever stopped working.
func TestCheckIntentPairedStaysQuiet(t *testing.T) {
	files := []string{"server/pkg/mcpkit/registry.go", "server/pkg/mcpkit/mcpkit.intent.md"}

	if got := checkIntentPaired(t.Context(), "/nonexistent", files); len(got) != 0 {
		t.Errorf("expected no warning when the intent change is already in the commit, got %v", got)
	}
}
