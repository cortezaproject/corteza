// Package tools holds the developer MCP's tools.
//
// The same split as the configurator side: a tool's declaration and its handler
// live together per subject file (branch.go, later test.go, intent.go), and
// Register is the single place that wires them. See ../SPEC.md for the families
// and CONVENTIONS.md under server/system/agentic/mcp for the authoring rules,
// which apply here too.
package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
)

// Register wires every developer tool into the registry.
//
// root is threaded through rather than read from the environment inside each
// tool, so that a test can point the whole surface at a scratch checkout.
func Register(reg *mcpkit.Registry, root string) {
	registerBranchStatus(reg, root)
	registerTestRun(reg, root)
	registerFormat(reg, root)
	registerIntentCheck(reg, root)
	registerIntentGoverning(reg, root)
	registerIntentAffected(reg, root)
}

// runner executes one command in the repository and returns its trimmed stdout.
type runner func(ctx context.Context, args ...string) (string, error)

// gitRunner binds git to a checkout.
//
// Every invocation is explicitly scoped with -C: inheriting the process working
// directory would make the answer depend on where the client happened to launch
// the server, which is precisely the ambiguity the root argument exists to
// remove.
func gitRunner(root string) runner {
	return func(ctx context.Context, args ...string) (string, error) {
		return run(ctx, root, "git", append([]string{"-C", root}, args...)...)
	}
}

// runAllowFail is run for a command whose non-zero exit is a result, not a
// failure.
//
// `intent check` exits non-zero exactly when it has drift to report, and `git
// diff --quiet` exits non-zero exactly when there is a diff. Using run for
// either throws away the answer and reports success, which is how the intent
// tool first came to say "no drift" against 51 drifted files.
func runAllowFail(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	// Both streams, because a CLI's choice between them is its own business and
	// a caller here wants what it said. `intent check` writes its entire report
	// to stderr, so reading stdout alone reports a clean tree over 51 drifted
	// files.
	var buf strings.Builder
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	return strings.TrimRight(buf.String(), "\r\n"), err
}

// run executes a command and returns its stdout.
//
// stderr is folded into the error rather than discarded: a coding agent reading
// "exit status 128" learns nothing, while "not a git repository" tells it what
// to do next.
//
// Only trailing newlines are trimmed, never leading whitespace: git's porcelain
// format encodes the index state in column one, so a leading space is data. A
// TrimSpace here silently shifted every field of every status line by one and
// reported unstaged files as staged.
func run(ctx context.Context, dir, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir

	var stderr strings.Builder
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), msg)
		}
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}

	return strings.TrimRight(string(out), "\r\n"), nil
}
