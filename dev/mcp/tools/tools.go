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
	"path/filepath"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// Register wires every developer tool into the registry.
//
// root is threaded through rather than read from the environment inside each
// tool, so that a test can point the whole surface at a scratch checkout. It is
// the default rather than the only answer: a tool acting on the repository
// takes a worktree argument that names another checkout of it.
func Register(reg *mcpkit.Registry, root string) {
	registerBranchStatus(reg, root)
	registerTestRun(reg, root)
	registerFormat(reg, root)
	registerIntentCheck(reg, root)
	registerIntentGoverning(reg, root)
	registerIntentAffected(reg, root)
	registerCommit(reg, root)
	registerServerStatus(reg, root)
	registerServerLogs(reg, root)
	registerFixtureCleanup(reg, root)
	registerUIVerify(reg, root)
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

// checkoutOption is the argument every tool that acts on a checkout declares.
//
// Declared once so the wording is identical across the surface: a caller
// reading two schemas must not have to work out whether they mean the same
// thing by it.
func checkoutOption() mcp.ToolOption {
	return mcp.WithString("worktree", mcp.Description(
		"Which checkout to act on: a worktree's directory name, its branch, or an absolute path. Omit "+
			"for the checkout this server was launched from, which is what a session working in the repo "+
			"root wants. Name it when the work lives in a git worktree and the session does not — an "+
			"orchestrated lane, or any task driving a branch checked out elsewhere: this server is one "+
			"process per session with no working directory of its own, so without a name every call here "+
			"acts on the launch checkout and a lane's commit lands on the wrong branch. Only checkouts of "+
			"this same repository resolve."))
}

// checkoutFor resolves the checkout one call acts on.
//
// Arguments are read leniently rather than through toolkit.Args, because a tool
// whose only parameter is this one is routinely called with no arguments object
// at all and refusing that would break every existing caller.
func checkoutFor(ctx context.Context, root string, req mcp.CallToolRequest) (string, error) {
	args, _ := req.Params.Arguments.(map[string]any)

	return resolveCheckout(ctx, root, toolkit.Str(args, "worktree"))
}

// checkout is one working tree of the repository.
type checkout struct {
	path   string
	branch string
}

// resolveCheckout maps a worktree argument to the checkout it names.
//
// The server is a stdio child launched once per client session from that
// session's working directory, so its root is fixed for the life of the process
// while the work a session does is not: a session sitting in the primary
// checkout routinely drives branches that live in worktrees. Naming the
// checkout per call is the only per-call answer available — a shared mutable
// "current checkout" would be read by every concurrent lane in the session, so
// one lane switching it would redirect another lane's commit.
//
// Resolution is closed over `git worktree list`, so a name can only ever reach a
// checkout of this same repository: pointing the commit tool at an unrelated
// directory is not something a caller can express.
func resolveCheckout(ctx context.Context, root, name string) (string, error) {
	if name == "" {
		return root, nil
	}

	list, err := listCheckouts(ctx, root)
	if err != nil {
		return "", err
	}

	// Tiered rather than one pass, so that a branch sharing a name with another
	// worktree's directory is resolved by precedence instead of refused as
	// ambiguous. Only a genuine tie within one tier is unanswerable.
	want := filepath.Clean(name)
	for _, match := range []func(checkout) bool{
		func(c checkout) bool { return filepath.Clean(c.path) == want },
		func(c checkout) bool { return c.branch == name },
		func(c checkout) bool { return filepath.Base(c.path) == name },
	} {
		var hits []string
		for _, c := range list {
			if match(c) {
				hits = append(hits, c.path)
			}
		}

		switch len(hits) {
		case 0:
			continue
		case 1:
			return hits[0], nil
		default:
			return "", fmt.Errorf("%q names %d checkouts (%s) — pass the path instead",
				name, len(hits), strings.Join(hits, ", "))
		}
	}

	return "", fmt.Errorf("no checkout of this repository is named %q; known: %s",
		name, strings.Join(checkoutNames(list), ", "))
}

// listCheckouts is every working tree git knows about, the primary first.
func listCheckouts(ctx context.Context, root string) ([]checkout, error) {
	out, err := gitRunner(root)(ctx, "worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}

	var (
		list []checkout
		cur  checkout
	)

	// Each record opens with its worktree line and closes at the next one; a
	// detached head simply has no branch line.
	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			if cur.path != "" {
				list = append(list, cur)
			}
			cur = checkout{path: strings.TrimPrefix(line, "worktree ")}
		case strings.HasPrefix(line, "branch "):
			cur.branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		}
	}

	if cur.path != "" {
		list = append(list, cur)
	}

	return list, nil
}

// checkoutNames is what a caller may write, one entry per checkout.
func checkoutNames(list []checkout) []string {
	out := make([]string, 0, len(list))

	for _, c := range list {
		name := filepath.Base(c.path)
		if c.branch != "" && c.branch != name {
			name += " (" + c.branch + ")"
		}
		out = append(out, name)
	}

	return out
}
