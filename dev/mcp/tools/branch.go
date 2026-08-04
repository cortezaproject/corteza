package tools

import (
	"context"
	"strconv"
	"strings"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/system/agentic/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// branchStatus is the shape dev_branch_status returns.
//
// Ahead/behind are counted against the main branch rather than the upstream:
// what a coding agent needs to know is how far the work has diverged from where
// it will merge, and a feature branch here often has no upstream at all.
type branchStatus struct {
	Branch     string   `json:"branch"`
	Base       string   `json:"base"`
	Ahead      int      `json:"ahead"`
	Behind     int      `json:"behind"`
	Staged     []string `json:"staged,omitempty"`
	Modified   []string `json:"modified,omitempty"`
	Untracked  []string `json:"untracked,omitempty"`
	Clean      bool     `json:"clean"`
	LastCommit string   `json:"lastCommit,omitempty"`
}

func registerBranchStatus(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_branch_status",
			mcp.WithDescription(
				"Where the working tree stands: current branch, how far it has diverged from the base branch "+
					"it will merge into, and what is staged, modified or untracked. Read this before committing "+
					"or reviewing, rather than running git status and reading the output.",
			),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Branch status",
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return branchStatusResult(ctx, root)
		},
	)
}

func branchStatusResult(ctx context.Context, root string) (*mcp.CallToolResult, error) {
	git := gitRunner(root)

	branch, err := git(ctx, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, err
	}

	out := branchStatus{Branch: branch, Base: baseBranch(ctx, git)}

	// A missing merge base is not an error: a freshly initialised repo, or a
	// branch with no common ancestor, simply has nothing to count.
	if counts, err := git(ctx, "rev-list", "--left-right", "--count", out.Base+"..."+out.Branch); err == nil {
		if fields := strings.Fields(counts); len(fields) == 2 {
			out.Behind, _ = strconv.Atoi(fields[0])
			out.Ahead, _ = strconv.Atoi(fields[1])
		}
	}

	status, err := git(ctx, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(status, "\n") {
		if len(line) < 4 {
			continue
		}

		// Porcelain v1: index status, worktree status, space, path. A path is
		// reported in both lists when it is staged and then modified again,
		// which is exactly the case a commit tool needs to see.
		x, y, path := line[0], line[1], strings.TrimSpace(line[3:])

		switch {
		case x == '?' && y == '?':
			out.Untracked = append(out.Untracked, path)
		default:
			if x != ' ' {
				out.Staged = append(out.Staged, path)
			}
			if y != ' ' {
				out.Modified = append(out.Modified, path)
			}
		}
	}

	out.Clean = len(out.Staged)+len(out.Modified)+len(out.Untracked) == 0

	if subject, err := git(ctx, "log", "-1", "--pretty=%h %s"); err == nil {
		out.LastCommit = subject
	}

	return toolkit.JSONResult(out)
}

// baseBranch is the branch this work merges into.
//
// The repo's main branch is a release line (2026.3.x today), not "main", so it
// cannot be hardcoded and origin/HEAD is often unset in a working clone. Asking
// git for the remote's default and falling back through the candidates is the
// only thing that stays correct as the release line moves.
func baseBranch(ctx context.Context, git runner) string {
	if ref, err := git(ctx, "rev-parse", "--abbrev-ref", "origin/HEAD"); err == nil && ref != "" {
		return ref
	}

	for _, candidate := range []string{"origin/main", "main", "master"} {
		if _, err := git(ctx, "rev-parse", "--verify", "--quiet", candidate); err == nil {
			return candidate
		}
	}

	return "HEAD"
}
