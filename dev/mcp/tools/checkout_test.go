package tools

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveCheckout pins which checkout a worktree argument reaches.
//
// The failure this prevents is silent and expensive: a session launched in the
// primary, driving a branch that lives in a worktree, gets the primary back
// from every tool here. Tests then run against the wrong tree and a commit
// lands on the wrong branch, both reporting success.
func TestResolveCheckout(t *testing.T) {
	base := gitFixture(t)

	primary := filepath.Join(base, "primary")
	lane := filepath.Join(base, "lane-x")
	named := filepath.Join(base, "dir-y")

	cases := []struct {
		name string
		arg  string
		want string // resolved path; empty means expect a refusal
		err  string // substring of the expected refusal
	}{
		{name: "omitted is the launch checkout", arg: "", want: primary},
		{name: "worktree directory name", arg: "lane-x", want: lane},
		{name: "worktree branch name", arg: "branch-y", want: named},
		{name: "worktree directory whose branch differs", arg: "dir-y", want: named},
		{name: "absolute path", arg: lane, want: lane},
		{name: "the primary by name", arg: "primary", want: primary},

		{name: "unknown name", arg: "no-such-lane", err: "no checkout of this repository"},
		{name: "a directory outside the repository", arg: t.TempDir(), err: "no checkout of this repository"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := resolveCheckout(context.Background(), primary, c.arg)

			if c.err != "" {
				if err == nil {
					t.Fatalf("resolveCheckout(%q) = %q, want a refusal naming %q", c.arg, got, c.err)
				}
				if !strings.Contains(err.Error(), c.err) {
					t.Fatalf("refusal %q does not name %q", err, c.err)
				}
				return
			}

			if err != nil {
				t.Fatalf("resolveCheckout(%q): %v", c.arg, err)
			}
			if got != c.want {
				t.Fatalf("resolveCheckout(%q) = %q, want %q", c.arg, got, c.want)
			}
		})
	}
}

// TestResolveCheckoutAmbiguous covers two checkouts sharing a directory name,
// which is only reachable from different parents. Answering one of them at
// random is the one outcome worse than refusing.
func TestResolveCheckoutAmbiguous(t *testing.T) {
	base := gitFixture(t)
	primary := filepath.Join(base, "primary")

	for _, p := range []struct{ path, branch string }{
		{filepath.Join(base, "a", "twin"), "twin-a"},
		{filepath.Join(base, "b", "twin"), "twin-b"},
	} {
		git(t, primary, "worktree", "add", "-b", p.branch, p.path)
	}

	if _, err := resolveCheckout(context.Background(), primary, "twin"); err == nil {
		t.Fatal("an ambiguous name resolved; want a refusal")
	} else if !strings.Contains(err.Error(), "names 2 checkouts") {
		t.Fatalf("refusal %q does not say which names collided", err)
	}
}

// gitFixture builds a repository with a primary checkout and two worktrees, one
// of them on a branch that does not match its directory name.
func gitFixture(t *testing.T) string {
	t.Helper()

	// The path git reports is fully resolved, so a symlinked temp dir would
	// make every comparison here fail for a reason that is not the subject.
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	primary := filepath.Join(base, "primary")

	git(t, base, "init", "-b", "main", primary)
	git(t, primary, "commit", "--allow-empty", "-m", "Root")
	git(t, primary, "worktree", "add", "-b", "lane-x", filepath.Join(base, "lane-x"))
	git(t, primary, "worktree", "add", "-b", "branch-y", filepath.Join(base, "dir-y"))

	return base
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()

	argv := append([]string{
		"-c", "user.name=Test",
		"-c", "user.email=test@example.com",
		"-C", dir,
	}, args...)

	if out, err := exec.Command("git", argv...).CombinedOutput(); err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
}
