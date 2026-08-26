package tools

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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
//
// `make watch` runs gin with --immediate, so it rebuilds AND respawns on any
// .go write on its own — no request to its proxy is involved, and a request to
// the proxy does not start a build. What it drops is an edit that lands while a
// build is already running: gin stamps its watch clock after the build returns,
// which puts that file's mtime in the past, and it is then skipped for good.
//
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

	WaitedSeconds int    `json:"waitedSeconds,omitempty"`
	Note          string `json:"note"`
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
			mcp.WithString("wait", mcp.Description(
				"Seconds to block for, as a string, while the running process still predates your edit. A "+
					"rebuild takes roughly 15s, so '60' covers one comfortably. Use it instead of polling: "+
					"this returns the moment the process serving requests started after every Go source, or "+
					"at the deadline with the reason it did not. Omit for a snapshot.")),
			mcpkit.InGroup(mcpkit.GroupDevelopment),
			mcpkit.WithRisk(mcpkit.RiskRead),
		),
		"Dev server status",
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args, err := toolkit.Args(req)
			if err != nil {
				return nil, err
			}

			wait, err := waitBudget(toolkit.Str(args, "wait"))
			if err != nil {
				return nil, err
			}

			return toolkit.JSONResult(awaitServer(ctx, root, wait))
		},
	)
}

// waitBudget reads the wait argument, which arrives as a string like every
// other numeric option on this surface.
//
// Capped rather than trusted: a caller that asks for an hour blocks the session
// on a server nobody is going to restart, and the useful answer — "it is still
// stale, here is why" — is available long before that.
func waitBudget(raw string) (time.Duration, error) {
	if raw == "" {
		return 0, nil
	}

	secs, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, fmt.Errorf("wait must be a number of seconds, got %q", raw)
	}

	if secs < 0 {
		return 0, fmt.Errorf("wait must not be negative, got %q", raw)
	}

	if secs > maxWaitSeconds {
		secs = maxWaitSeconds
	}

	return time.Duration(secs) * time.Second, nil
}

const maxWaitSeconds = 180

// How long to leave between checks while waiting. A variable so the wait loop
// can be exercised without a test that sleeps for the length of a real build.
var waitPoll = 2 * time.Second

// awaitServer polls checkServer until the server is up and running the caller's
// code, or the budget runs out.
//
// The alternative is what agents do without it: call the tool, read "stale",
// sleep by hand, call again. That loop is written differently every time and
// usually gives up one poll early, which reports a fix as not working.
func awaitServer(ctx context.Context, root string, budget time.Duration) serverStatus {
	started := time.Now()

	for {
		out := checkServer(ctx, root)

		if budget == 0 {
			return out
		}

		waited := time.Since(started)
		out.WaitedSeconds = int(waited.Round(time.Second) / time.Second)

		if out.Up && !out.Stale {
			return out
		}

		if waited+waitPoll > budget {
			out.Note = "waited " + strconv.Itoa(out.WaitedSeconds) + "s and " + out.Note
			return out
		}

		select {
		case <-ctx.Done():
			return out
		case <-time.After(waitPoll):
		}
	}
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

	out.Note = statusNote(out, newest, haveProcess, watcherRunning(ctx, root))

	return out
}

// statusNote says what the timestamps mean and what to do about them.
//
// It is the widest-read description of gin's behaviour in the repo — a session
// calls this tool, reads the note, and acts on it — so what it says has to be
// what gin actually does. A note prescribing a poke at gin's proxy port sends
// every session curling it for a build that has already happened without them,
// and they carry the instruction onwards as fact.
func statusNote(out serverStatus, newest string, haveProcess, watcherRunning bool) string {
	switch {
	case !out.Up && watcherRunning:
		return "the dev server is not answering, but gin IS running — it is most likely mid-restart, so " +
			"check again before concluding anything, or call this with wait. Do NOT start a second " +
			"'make watch': it will fail with 'address already in use' and its 'tee build/dev.log' " +
			"truncates the log on the way out"
	case !out.Up:
		return "the dev server is not answering and no gin watcher is running. Start it with " +
			"'cd server && make watch', or ask the human to"
	case out.Stale && haveProcess:
		return fmt.Sprintf("the process serving requests started BEFORE %s was last edited, so it is "+
			"running code older than your change. Anything you reproduce against it right now describes the "+
			"OLD build — including a failure, which will look exactly like a fix that did not work. gin "+
			"rebuilds and respawns on its own, so the first move is to wait about 15s (call this with "+
			"wait). Poking gin's proxy port does nothing here and is not the remedy. What gin drops is an "+
			"edit that landed WHILE a build was running: it is skipped for good, and only re-saving the "+
			"file gets it built. That restarts the server, so on the shared primary it is the human's call. "+
			"If nothing comes back, read server/build/dev.log with dev_server_logs — a failed build leaves "+
			"the old process serving happily", newest)
	case out.Stale:
		return fmt.Sprintf("the binary is OLDER than %s and no running process was found to check against. "+
			"It is still serving the previous build — check server/build/dev.log with dev_server_logs, "+
			"because a failed rebuild leaves the old binary running", newest)
	case haveProcess:
		return "up, and the process serving requests started after every Go source under server/ was last " +
			"changed — this is the check that actually means it is running your code"
	default:
		return "up, and the binary is newer than every Go source under server/. No running process was " +
			"found to confirm against, so this is the weaker of the two checks"
	}
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
	pids, err := runAllowFail(ctx, root, "pgrep", "-f", serverProcessPattern(root))
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

// serverProcessPattern matches THIS checkout's server child and nothing else.
//
// Matching the process name (`pgrep -x gin-bin`) matched any checkout's: with a
// second Corteza tree also running one, this reported that server's start time
// and called a Human server restarted two hours later stale, while the API was
// already serving the fix. gin runs the child by absolute path, so anchoring on
// the path separates them — and separates the child from the watcher and from
// the shell wrapping the pipeline, both of which carry a relative
// "build/gin-bin" in their command line.
func serverProcessPattern(root string) string {
	bin := filepath.Join(root, "server", "build", "gin-bin")

	return "^" + regexp.QuoteMeta(bin) + "( |$)"
}

// watcherRunning answers whether THIS checkout has a gin watcher up, which is
// the difference between "mid-restart, ask again" and "nothing is running".
//
// Scoped by the port gin was told to proxy to, which worktree.sh assigns per
// slot: a bare `pgrep -x gin` finds any checkout's watcher, and would tell a
// lane whose server is down that it is merely restarting because slot 0's is up.
func watcherRunning(ctx context.Context, root string) bool {
	out, err := runAllowFail(ctx, root, "pgrep", "-f", watcherPattern(devServerPort(root)))
	return err == nil && strings.TrimSpace(out) != ""
}

func watcherPattern(port string) string {
	return "gin .*--appPort " + regexp.QuoteMeta(port) + "( |$)"
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
	return "http://localhost:" + devServerPort(root)
}

func devServerPort(root string) string {
	// The port a checkout's server listens on is the one in its own server/.env
	// — worktree.sh writes it per slot, so a lane answers for its own server
	// rather than the primary's. Falling back rather than failing keeps this
	// useful on a checkout that has not been set up yet.
	// The LAST assignment wins, because that is the one the server gets:
	// godotenv parses the whole file into a map before applying it. Stopping at
	// the first match names a port nothing is bound to as soon as somebody
	// appends an override at the end of the file.
	found := ""

	if raw, err := os.ReadFile(filepath.Join(root, "server", ".env")); err == nil {
		for _, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "HTTP_ADDR=") {
				continue
			}

			// HTTP_ADDR is a bind address: ":1543", "127.0.0.1:1543". An
			// empty assignment has no fields at all, so do not index blindly.
			fields := strings.Fields(strings.TrimPrefix(line, "HTTP_ADDR="))
			if len(fields) == 0 {
				continue
			}

			if _, port, ok := strings.Cut(strings.Trim(fields[0], `"'`), ":"); ok && port != "" {
				found = port
			}
		}
	}

	if found != "" {
		return found
	}

	return "1043"
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
