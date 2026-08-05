package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/crusttech/human/server/pkg/mcpkit"
	"github.com/crusttech/human/server/pkg/mcpkit/toolkit"
	"github.com/mark3labs/mcp-go/mcp"
)

// serverStatus answers the question that has to come before any reproduction:
// is the thing I am about to test actually running my code?
//
// A reproduction against a stale binary is worse than no reproduction, because
// it looks like evidence. `gin` rebuilds on change, but a build that failed
// leaves the previous binary serving happily, and nothing about the API's
// responses says so.
type serverStatus struct {
	Up      bool   `json:"up"`
	URL     string `json:"url"`
	Healthy bool   `json:"healthy"`

	BinaryBuiltAt string `json:"binaryBuiltAt,omitempty"`
	NewestSource  string `json:"newestSource,omitempty"`
	Stale         bool   `json:"stale"`

	Note string `json:"note"`
}

func registerServerStatus(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_server_status",
			mcp.WithDescription(
				"Is the local Human dev server up, and is it running your code? Call this before reproducing "+
					"anything against it. A stale binary answers requests normally while serving code you have "+
					"already changed, so a reproduction against one proves nothing and looks like it proves "+
					"something. If this reports stale, wait for the rebuild rather than drawing conclusions.",
			),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Dev server status",
		func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return toolkit.JSONResult(checkServer(ctx, root))
		},
	)
}

func checkServer(ctx context.Context, root string) serverStatus {
	out := serverStatus{URL: devServerURL(root)}

	// curl rather than net/http: the point is to ask the way every other tool
	// in dev/agent asks, so a proxy or a host override applies here too.
	if _, err := run(ctx, root, "curl", "-sf", "-m", "5", "-o", "/dev/null", out.URL+"/healthcheck"); err == nil {
		out.Up = true
		out.Healthy = true
	}

	bin := filepath.Join(root, "server", "build", "gin-bin")
	info, err := os.Stat(bin)
	if err != nil {
		out.Note = "no dev binary at server/build/gin-bin — the dev server is started with 'cd server && make watch'"
		return out
	}

	out.BinaryBuiltAt = info.ModTime().Format(time.RFC3339)

	// Newest Go source under server/, excluding vendor: if anything is newer
	// than the binary, gin has not finished rebuilding — or the rebuild failed.
	newest, at := newestGoSource(ctx, root)
	if newest != "" {
		out.NewestSource = newest
		out.Stale = at.After(info.ModTime())
	}

	switch {
	case !out.Up:
		out.Note = "the dev server is not answering. Start it with 'cd server && make watch', or ask the human to"
	case out.Stale:
		out.Note = fmt.Sprintf("the running binary is OLDER than %s. It is still serving the previous build — "+
			"wait for gin to finish, and check server/build/dev.log with dev_server_logs if it does not, "+
			"because a failed rebuild leaves the old binary running", newest)
	default:
		out.Note = "up, and the binary is newer than every Go source under server/"
	}

	return out
}

func newestGoSource(ctx context.Context, root string) (string, time.Time) {
	out, err := run(ctx, root, "find", filepath.Join(root, "server"),
		"-name", "*.go", "-not", "-path", "*/vendor/*", "-newer",
		filepath.Join(root, "server", "build", "gin-bin"), "-print")
	if err != nil || strings.TrimSpace(out) == "" {
		return "", time.Time{}
	}

	newest := strings.Split(strings.TrimSpace(out), "\n")[0]
	info, err := os.Stat(newest)
	if err != nil {
		return "", time.Time{}
	}

	return strings.TrimPrefix(strings.TrimPrefix(newest, root), "/"), info.ModTime()
}

func devServerURL(root string) string {
	// The port lives in the dev toolkit's shared config; falling back rather
	// than failing keeps this useful on a checkout that has not bootstrapped.
	if raw, err := os.ReadFile(filepath.Join(root, "dev", "agent", "common.sh")); err == nil {
		for _, line := range strings.Split(string(raw), "\n") {
			if _, value, ok := strings.Cut(line, "HUMAN_BASE_URL="); ok {
				if v := strings.Trim(strings.Fields(value)[0], `"'`); strings.HasPrefix(v, "http") {
					return v
				}
			}
		}
	}

	return "http://localhost:1043"
}

func registerServerLogs(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_server_logs",
			mcp.WithDescription(
				"Tail the dev server log, optionally filtered. Use it when a reproduction fails in a way the "+
					"API response does not explain — a panic, a failed rebuild, a store error swallowed on the "+
					"way out. The server logs far more than it returns.",
			),
			mcp.WithString("pattern", mcp.Description(
				"Case-insensitive filter. Omit for the plain tail.")),
			mcp.WithString("lines", mcp.Description(
				"How many lines to read, as a string. Defaults to 100; the result is capped regardless.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Dev server logs",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			argv := []string{"agent/logs.sh"}
			if n := toolkit.Str(args, "lines"); n != "" {
				if _, err := strconv.Atoi(n); err != nil {
					return nil, fmt.Errorf("lines must be a number, got %q", n)
				}
				argv = append(argv, "-n", n)
			}
			if p := toolkit.Str(args, "pattern"); p != "" {
				argv = append(argv, p)
			}

			out, err := runAllowFail(ctx, filepath.Join(root, "dev"), "bash", argv...)
			if strings.TrimSpace(out) == "" {
				if err != nil {
					return nil, fmt.Errorf("cannot read the dev log: %w", err)
				}
				return toolkit.JSONResult(map[string]any{"lines": []string{}, "note": "nothing matched"})
			}

			return toolkit.JSONResult(map[string]any{
				"lines": nonEmptyLines(firstLines(out, 200)),
				"note":  "tail of server/build/dev.log, newest last",
			})
		},
	)
}

type cleanupReport struct {
	Removed []string `json:"removed"`
	Note    string   `json:"note"`
}

func registerFixtureCleanup(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_fixture_cleanup",
			mcp.WithDescription(
				"Delete everything you created on the dev server, which means every compose namespace whose "+
					"slug starts with 'agent-'. Unprefixed data is never touched — it is not yours. Run this "+
					"when a task is done: state left behind pollutes the next reproduction and the next "+
					"lookup, and nothing else removes it.",
			),
			mcp.WithBoolean("purge", mcp.Description(
				"Also hard-delete soft-deleted agent- namespaces. Repeated cycles leave slug-sharing corpses "+
					"that make lookups confusing. Off by default because it is irreversible.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskDestructive),
		),
		"Clean up agent data",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			argv := []string{"agent/cleanup.sh"}
			if toolkit.Bool(args, "purge") {
				argv = append(argv, "--purge")
			}

			out, err := runAllowFail(ctx, filepath.Join(root, "dev"), "bash", argv...)

			// Only lines in the script's own format count as removals.
			//
			// Handing back whatever the script printed was wrong in the exact
			// way this package keeps finding elsewhere: when the dev server was
			// mid-restart the script's curl returned nothing, its Python threw,
			// and the traceback came back to the caller as a list of things
			// that had supposedly been deleted.
			report := cleanupReport{Removed: []string{}}
			var noise []string

			for _, line := range nonEmptyLines(out) {
				switch {
				case strings.HasPrefix(line, "deleted "), strings.HasPrefix(line, "purged "):
					report.Removed = append(report.Removed, strings.TrimSpace(line))
				case line == "cleanup done":
				default:
					noise = append(noise, line)
				}
			}

			if err != nil || len(noise) > 0 {
				return nil, fmt.Errorf(
					"cleanup did not complete — the dev server may be restarting, check dev_server_status: %s",
					firstLines(strings.Join(noise, "\n"), 5),
				)
			}

			if len(report.Removed) == 0 {
				report.Note = "nothing to remove — no agent- prefixed data on the dev server"
			} else {
				report.Note = fmt.Sprintf("removed %d item(s); unprefixed data was not touched", len(report.Removed))
			}

			return toolkit.JSONResult(report)
		},
	)
}
