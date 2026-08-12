package tools

import (
	"context"
	"fmt"
	"io/fs"
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
// The authority on "is it running my code" is the START TIME OF THE PROCESS
// serving requests, never a file mtime. Comparing the binary's mtime to the
// newest source looks equivalent and is not: gin can begin a rebuild before an
// edit lands and finish after it, which leaves a binary that is newer than every
// source while containing the older code. That reported "up, and the binary is
// newer than every Go source" for a full cycle while the API kept returning
// pre-fix results, and the re-run of a correct fix looked like a failed fix.
type serverStatus struct {
	Up      bool   `json:"up"`
	URL     string `json:"url"`
	Healthy bool   `json:"healthy"`

	BinaryBuiltAt    string `json:"binaryBuiltAt,omitempty"`
	ProcessStartedAt string `json:"processStartedAt,omitempty"`
	NewestSource     string `json:"newestSource,omitempty"`
	NewestSourceAt   string `json:"newestSourceAt,omitempty"`
	Stale            bool   `json:"stale"`

	Note string `json:"note"`
}

func registerServerStatus(reg *mcpkit.Registry, root string) {
	reg.RegisterTool(
		mcp.NewTool("dev_server_status",
			mcp.WithDescription(
				"Is the local Human dev server up, and is it running your code? Call this before reproducing "+
					"anything against it, AND again before believing a re-run of the same reproduction after a "+
					"server-side fix. A stale process answers requests normally while serving code you have "+
					"already changed, so a reproduction against one proves nothing and looks like it proves "+
					"something — in both directions: the failure you still see may be the old build rather "+
					"than a fix that did not work. Staleness here is decided by when the running process "+
					"started, not by file mtimes, because a rebuild that begins before an edit and ends after "+
					"it leaves a binary newer than every source while running the older code.",
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

	// Newest Go source under server/, excluding vendor.
	newest, at := newestGoSource(root)
	if newest != "" {
		out.NewestSource = newest
		out.NewestSourceAt = at.Format(time.RFC3339)
	}

	started, haveProcess := serverProcessStart(ctx, root)
	if haveProcess {
		out.ProcessStartedAt = started.Format(time.RFC3339)
	}

	out.Stale = staleAgainst(at, started, info.ModTime(), haveProcess)

	watcherRunning := processRunning(ctx, root, "gin")

	switch {
	case !out.Up && watcherRunning:
		out.Note = "the dev server is not answering, but gin IS running — it is most likely mid-restart, so " +
			"check again before concluding anything. Do NOT start a second 'make watch': it will fail with " +
			"'address already in use' and its 'tee build/dev.log' truncates the log on the way out. gin " +
			"rebuilds lazily when its proxy is hit, so 'curl -s localhost:3001/api/' prods it along"
	case !out.Up:
		out.Note = "the dev server is not answering and no gin watcher is running. Start it with " +
			"'cd server && make watch', or ask the human to"
	case out.Stale && haveProcess:
		out.Note = fmt.Sprintf("the process serving requests started BEFORE %s was last edited, so it is "+
			"running code older than your change. Anything you reproduce against it right now describes the "+
			"OLD build — including a failure, which will look exactly like a fix that did not work. gin "+
			"rebuilds lazily when its proxy is hit ('curl -s localhost:3001/api/'); requests to the API port "+
			"never trigger it. If it does not come back, read server/build/dev.log with dev_server_logs, "+
			"because a failed rebuild leaves the old process serving happily", newest)
	case out.Stale:
		out.Note = fmt.Sprintf("the binary is OLDER than %s and no running process was found to check against. "+
			"It is still serving the previous build — check server/build/dev.log with dev_server_logs, "+
			"because a failed rebuild leaves the old binary running", newest)
	case haveProcess:
		out.Note = "up, and the process serving requests started after every Go source under server/ was last " +
			"changed — this is the check that actually means it is running your code"
	default:
		out.Note = "up, and the binary is newer than every Go source under server/. No running process was " +
			"found to confirm against, so this is the weaker of the two checks"
	}

	return out
}

// staleAgainst decides whether what is answering requests predates the newest
// source, preferring the running process over the binary on disk.
//
// The two are not interchangeable. gin can start compiling before an edit lands
// and finish after it, which produces a binary whose mtime is newer than every
// source while its contents are older — so the mtime comparison reports fresh
// for a process that is serving the previous build. It is only the fallback.
func staleAgainst(newestSourceAt, processStartedAt, binaryBuiltAt time.Time, haveProcess bool) bool {
	if newestSourceAt.IsZero() {
		return false
	}

	if haveProcess {
		return newestSourceAt.After(processStartedAt)
	}

	return newestSourceAt.After(binaryBuiltAt)
}

// serverProcessStart returns when the process actually answering requests
// started. That is the only thing that settles whether a reproduction is
// testing your code, and it is not derivable from any file's mtime.
//
// Elapsed seconds rather than ps's start-time column: the latter is formatted
// per locale and would need parsing, while etimes is an integer everywhere.
func serverProcessStart(ctx context.Context, root string) (time.Time, bool) {
	// -x matches the process name exactly, which separates the child (gin-bin)
	// from the watcher (gin) and from the shell wrapping the pipeline, all three
	// of which carry "build/gin-bin" somewhere in their command line.
	pids, err := runAllowFail(ctx, root, "pgrep", "-x", "gin-bin")
	if err != nil || strings.TrimSpace(pids) == "" {
		return time.Time{}, false
	}

	fields := strings.Fields(pids)
	if len(fields) == 0 {
		return time.Time{}, false
	}

	elapsed, err := runAllowFail(ctx, root, "ps", "-o", "etimes=", "-p", fields[0])
	if err != nil {
		return time.Time{}, false
	}

	secs, err := strconv.Atoi(strings.TrimSpace(elapsed))
	if err != nil {
		return time.Time{}, false
	}

	return time.Now().Add(-time.Duration(secs) * time.Second), true
}

func processRunning(ctx context.Context, root, name string) bool {
	out, err := runAllowFail(ctx, root, "pgrep", "-x", name)
	return err == nil && strings.TrimSpace(out) != ""
}

// newestGoSource finds the most recently modified Go source under server/.
//
// Absolute rather than relative to the binary: the caller now compares it
// against the running process too, and "newer than the binary" cannot answer
// that. Walked in Go rather than shelled out to find, which keeps the mtime
// comparison in one place and needs no find flags beyond POSIX.
func newestGoSource(root string) (string, time.Time) {
	var (
		newestPath string
		newestAt   time.Time
	)

	_ = filepath.WalkDir(filepath.Join(root, "server"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			switch d.Name() {
			case "vendor", "build", "node_modules", ".git", "testdata":
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if info.ModTime().After(newestAt) {
			newestAt, newestPath = info.ModTime(), path
		}

		return nil
	})

	if newestPath == "" {
		return "", time.Time{}
	}

	return strings.TrimPrefix(strings.TrimPrefix(newestPath, root), "/"), newestAt
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
				"Delete what THIS SESSION created on the dev server. The candidates come from the ledger "+
					"dev/agent/.state/created.jsonl, which api.sh and mcp.py append to as they create things "+
					"— not from slug prefixes. Naming something 'agent-' neither adds it to the ledger nor "+
					"protects anything else: a namespace this session did not create is never a candidate, "+
					"whoever made it and whatever it is called, and one it did create is removed whatever it "+
					"is named. So there is no need to prefix what you create. Run this when a task is done: "+
					"state left behind pollutes the next reproduction and the next lookup.",
			),
			mcp.WithBoolean("purge", mcp.Description(
				"Also hard-delete soft-deleted namespaces. Repeated cycles leave slug-sharing corpses that "+
					"make lookups confusing. This part is NOT ledger-scoped — it hits every soft-deleted "+
					"namespace on the server, including ones someone else deleted. Off by default because it "+
					"is irreversible.")),
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
				report.Note = "nothing to remove — this session's ledger records nothing it created"
			} else {
				report.Note = fmt.Sprintf("removed %d item(s); unprefixed data was not touched", len(report.Removed))
			}

			return toolkit.JSONResult(report)
		},
	)
}
