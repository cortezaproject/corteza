// Command dev-mcp is the developer MCP: the L1 layer of Human's MCP system.
//
// It standardises how AI-assisted work on this repository is done — investigate,
// verify, review, commit, and exercise a running Human — by giving a coding
// agent typed primitives instead of raw shell output. See SPEC.md alongside this
// file for what belongs here and what stays a skill.
//
// It is a separate module, and a separate binary, on purpose. Its tools shell
// out to git, go and prettier, and that capability must not exist inside the
// binary customers run. A build tag would be a weaker guarantee: it has to be
// correct in every build path, whereas a module nothing imports cannot be linked
// in by accident. It also has to work with the Human server down — investigating
// and reviewing need no API.
//
// The machinery is shared with the configurator MCP through
// server/pkg/mcpkit, which may not import anything from Human's domain.
package main

import (
	"fmt"
	"os"

	"github.com/crusttech/human/dev/mcp/tools"
	"github.com/crusttech/human/server/pkg/mcpkit"
)

const (
	name    = "human-dev"
	version = "0.1.0"
)

func main() {
	// Repo root is passed in rather than discovered, because the client decides
	// the working directory it launches us from and guessing wrong would mean
	// running git and go against the wrong tree.
	root, err := repoRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dev-mcp:", err)
		os.Exit(1)
	}

	reg := mcpkit.NewRegistry()
	tools.Register(reg, root)

	// Stdio, not HTTP: the client launches this as a child process. Anything
	// written to stdout that is not protocol is a parse error at the other end,
	// so every diagnostic in this binary goes to stderr.
	if err := mcpkit.NewMCPServer(reg, name, version).ServeStdio(); err != nil {
		fmt.Fprintln(os.Stderr, "dev-mcp:", err)
		os.Exit(1)
	}
}

// repoRoot resolves the checkout this server operates on.
//
// HUMAN_REPO_ROOT wins so a client can be explicit; otherwise the working
// directory is used, which is what a client launching us from the repo gives.
func repoRoot() (string, error) {
	if root := os.Getenv("HUMAN_REPO_ROOT"); root != "" {
		return root, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine the repository root: %w (set HUMAN_REPO_ROOT)", err)
	}

	return wd, nil
}
