package tools

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// TestStaleAgainstPrefersTheRunningProcess pins the case that made this check
// worth having.
//
// A fix was on disk and its unit test was green, gin rebuilt, and the tool said
// "up, and the binary is newer than every Go source" — while the API kept
// returning pre-fix results. gin had begun compiling before the edit landed and
// finished after it, so the binary's mtime was newer than every source and its
// contents were older. Judged by mtime alone the server looked fresh; judged by
// when the process started, it plainly was not.
//
// The cost of getting this wrong is asymmetric: reporting fresh for a stale
// server turns the re-run of a correct fix into a red result, which reads as a
// fix that did not work.
func TestStaleAgainstPrefersTheRunningProcess(t *testing.T) {
	var (
		processStart = time.Date(2026, 8, 12, 14, 27, 49, 0, time.UTC)
		sourceEdit   = processStart.Add(2 * time.Minute)
		binaryBuilt  = sourceEdit.Add(30 * time.Second)
	)

	cases := []struct {
		name        string
		sourceAt    time.Time
		processAt   time.Time
		binaryAt    time.Time
		haveProcess bool
		want        bool
	}{
		{
			// mtimes say fresh, the process says otherwise.
			name:        "binary newer than source but process predates the edit",
			sourceAt:    sourceEdit,
			processAt:   processStart,
			binaryAt:    binaryBuilt,
			haveProcess: true,
			want:        true,
		},
		{
			name:        "process restarted after the edit",
			sourceAt:    sourceEdit,
			processAt:   binaryBuilt,
			binaryAt:    binaryBuilt,
			haveProcess: true,
			want:        false,
		},
		{
			// Without a process to ask, the binary is all there is.
			name:        "no process, binary older than source",
			sourceAt:    sourceEdit,
			processAt:   time.Time{},
			binaryAt:    processStart,
			haveProcess: false,
			want:        true,
		},
		{
			name:        "no process, binary newer than source",
			sourceAt:    sourceEdit,
			processAt:   time.Time{},
			binaryAt:    binaryBuilt,
			haveProcess: false,
			want:        false,
		},
		{
			// Nothing found to compare against is not evidence of staleness.
			name:        "no source timestamp at all",
			sourceAt:    time.Time{},
			processAt:   processStart,
			binaryAt:    binaryBuilt,
			haveProcess: true,
			want:        false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := staleAgainst(c.sourceAt, c.processAt, c.binaryAt, c.haveProcess); got != c.want {
				t.Fatalf("staleAgainst() = %v, want %v", got, c.want)
			}
		})
	}
}

// TestNewestGoSourceIgnoresTheBuildDirectory guards the walk against reporting
// something that is not a source file as the newest source, which would make
// every status permanently stale.
func TestNewestGoSourceIgnoresTheBuildDirectory(t *testing.T) {
	root := t.TempDir()

	writeAt := func(rel string, at time.Time) {
		t.Helper()

		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", rel, err)
		}
		if err := os.WriteFile(full, []byte("package x\n"), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
		if err := os.Chtimes(full, at, at); err != nil {
			t.Fatalf("chtimes %s: %v", rel, err)
		}
	}

	base := time.Date(2026, 8, 12, 12, 0, 0, 0, time.UTC)

	writeAt("server/compose/service/record.go", base)
	writeAt("server/pkg/dal/aggregate.go", base.Add(time.Hour))
	// Both of these are newer, and neither is a source file this cares about.
	writeAt("server/vendor/example.com/lib/lib.go", base.Add(2*time.Hour))
	writeAt("server/build/generated.go", base.Add(3*time.Hour))

	name, at := newestGoSource(root)

	if name != "server/pkg/dal/aggregate.go" {
		t.Fatalf("newest source = %q, want server/pkg/dal/aggregate.go", name)
	}

	if !at.Equal(base.Add(time.Hour)) {
		t.Fatalf("newest source time = %v, want %v", at, base.Add(time.Hour))
	}
}

// TestDevServerURLAnswersForItsOwnCheckout pins what this has to get right.
//
// It used to parse HUMAN_BASE_URL out of dev/agent/common.sh — a variable that
// is set nowhere in the repo — so every checkout fell through to the hardcoded
// primary port. A lane asking whether its server was up was told about slot 0's,
// and a stopped lane reported as running.
func TestDevServerURLAnswersForItsOwnCheckout(t *testing.T) {
	cases := []struct {
		name string
		env  string
		want string
	}{
		{
			name: "a worktree's own port, not the primary's",
			env:  "DOMAIN=localhost:1543\nHTTP_ADDR=:1543\nENVIRONMENT=dev\n",
			want: "http://localhost:1543",
		},
		{
			name: "a host-qualified bind address",
			env:  "HTTP_ADDR=127.0.0.1:1743\n",
			want: "http://localhost:1743",
		},
		{
			name: "a commented assignment is not the value",
			env:  "#HTTP_ADDR=:9999\nHTTP_ADDR=:1343\n",
			want: "http://localhost:1343",
		},
		{
			// godotenv parses the whole file into a map before applying it, so
			// the server listens on the last assignment, not the first.
			name: "the last assignment wins, as the server sees it",
			env:  "HTTP_ADDR=:1343\nHTTP_ADDR=:1443\n",
			want: "http://localhost:1443",
		},
		{
			name: "an empty assignment falls back rather than panicking",
			env:  "HTTP_ADDR=\n",
			want: "http://localhost:1043",
		},
		{
			name: "no env file at all",
			env:  "",
			want: "http://localhost:1043",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			root := t.TempDir()

			if c.env != "" {
				if err := os.MkdirAll(filepath.Join(root, "server"), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(root, "server", ".env"), []byte(c.env), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			if got := devServerURL(root); got != c.want {
				t.Errorf("devServerURL() = %q, want %q", got, c.want)
			}
		})
	}
}

// TestServerProcessPatternMatchesOnlyThisCheckout pins the mismatch that made
// the freshness flag lie.
//
// With a second Corteza tree also running a `gin-bin`, matching on the process
// name reported THAT process's start time and called a Human server restarted
// two hours later stale, while the API was already serving the fix.
func TestServerProcessPatternMatchesOnlyThisCheckout(t *testing.T) {
	pattern := regexp.MustCompile(serverProcessPattern("/home/dev/Human/human"))

	cases := []struct {
		name    string
		cmdline string
		want    bool
	}{
		{
			name:    "this checkout's child",
			cmdline: "/home/dev/Human/human/server/build/gin-bin --env-file .env serve",
			want:    true,
		},
		{
			name:    "another checkout's child",
			cmdline: "/home/dev/Corteza/server/build/gin-bin --env-file .env serve",
			want:    false,
		},
		{
			// A worktree is its own checkout with its own server and its own
			// port; the primary must not answer for it, or a lane's status is
			// slot 0's.
			name:    "a worktree's child",
			cmdline: "/home/dev/Human/human-lanes/fix-chart/server/build/gin-bin --env-file .env serve",
			want:    false,
		},
		{
			name:    "the watcher, which runs the binary by a relative path",
			cmdline: "/home/dev/go/bin/gin --laddr localhost --port 3001 --build cmd/human --immediate --bin build/gin-bin -- --env-file .env serve",
			want:    false,
		},
		{
			name:    "the shell wrapping the watcher's pipeline",
			cmdline: "/bin/sh -c /home/dev/go/bin/gin --bin build/gin-bin -- serve 2>&1 | tee -a build/dev.log",
			want:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := pattern.MatchString(c.cmdline); got != c.want {
				t.Fatalf("match(%q) = %v, want %v", c.cmdline, got, c.want)
			}
		})
	}
}

// TestWaitBudgetReadsSecondsAndCaps guards the one argument that can block a
// session: a wait nobody bounded is a session hung on a server nobody is going
// to restart.
func TestWaitBudgetReadsSecondsAndCaps(t *testing.T) {
	cases := []struct {
		raw     string
		want    time.Duration
		wantErr bool
	}{
		{raw: "", want: 0},
		{raw: "60", want: 60 * time.Second},
		{raw: " 15 ", want: 15 * time.Second},
		{raw: "0", want: 0},
		{raw: "99999", want: maxWaitSeconds * time.Second},
		{raw: "-5", wantErr: true},
		{raw: "a while", wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.raw, func(t *testing.T) {
			got, err := waitBudget(c.raw)

			if c.wantErr {
				if err == nil {
					t.Fatalf("waitBudget(%q) = %v, want an error", c.raw, got)
				}
				return
			}

			if err != nil {
				t.Fatalf("waitBudget(%q): %v", c.raw, err)
			}

			if got != c.want {
				t.Fatalf("waitBudget(%q) = %v, want %v", c.raw, got, c.want)
			}
		})
	}
}

// TestStatusNoteNeverPrescribesThePoke keeps the tool from teaching a remedy
// gin has not needed since --immediate.
//
// This note is the widest-read description of the watcher in the repo. While it
// said "gin rebuilds lazily when its proxy is hit", sessions curled the proxy
// port for a build that had already happened without them, and reported the
// instruction onwards as fact.
func TestStatusNoteNeverPrescribesThePoke(t *testing.T) {
	banned := []string{"lazily", "proxy is hit", "localhost:3001", "prods it"}

	states := []struct {
		name           string
		out            serverStatus
		haveProcess    bool
		watcherRunning bool
	}{
		{name: "down, watcher running", watcherRunning: true},
		{name: "down, nothing running"},
		{name: "up but stale", out: serverStatus{Up: true, Stale: true}, haveProcess: true},
		{name: "stale with no process", out: serverStatus{Up: true, Stale: true}},
		{name: "fresh", out: serverStatus{Up: true}, haveProcess: true},
		{name: "fresh, no process"},
	}

	for _, s := range states {
		t.Run(s.name, func(t *testing.T) {
			note := statusNote(s.out, "server/compose/service/record.go", s.haveProcess, s.watcherRunning)

			if note == "" {
				t.Fatal("every state must say something")
			}

			for _, bad := range banned {
				if strings.Contains(note, bad) {
					t.Errorf("note prescribes the poke (%q): %s", bad, note)
				}
			}
		})
	}
}

// TestAwaitServerPollsUntilTheBudgetRunsOut pins the behaviour that replaces a
// hand-rolled poll loop: keep checking, and when the answer is still no, say
// how long it waited rather than reporting a plain snapshot.
//
// Written against a checkout with nothing running, which is the shape of the
// case that matters — a session waiting on a rebuild that is never going to
// arrive has to be told so, not left holding a status it might read as fresh.
func TestAwaitServerPollsUntilTheBudgetRunsOut(t *testing.T) {
	restore := waitPoll
	waitPoll = 10 * time.Millisecond
	t.Cleanup(func() { waitPoll = restore })

	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "server", "build"), 0o755); err != nil {
		t.Fatal(err)
	}
	// A binary must exist, or the check returns before it reaches staleness.
	if err := os.WriteFile(filepath.Join(root, "server", "build", "gin-bin"), []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Nothing is listening on this port, so the server never comes up.
	if err := os.WriteFile(filepath.Join(root, "server", ".env"), []byte("HTTP_ADDR=:1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := awaitServer(context.Background(), root, 60*time.Millisecond)

	if out.Up {
		t.Fatal("nothing is listening, so it must not report up")
	}

	if !strings.HasPrefix(out.Note, "waited ") {
		t.Errorf("a timed-out wait must say it waited, got: %s", out.Note)
	}

	// Without the budget it returns the first check and never polls at all.
	snapshot := awaitServer(context.Background(), root, 0)

	if strings.HasPrefix(snapshot.Note, "waited ") {
		t.Errorf("a snapshot must not claim to have waited, got: %s", snapshot.Note)
	}
}

// TestWatcherPatternMatchesOnlyThisCheckoutsWatcher keeps a lane from being told
// its dead server is "mid-restart" because another slot's watcher is up.
//
// gin's proxy-to port is the one thing on its command line that worktree.sh
// assigns per slot, so it is what separates the watchers.
func TestWatcherPatternMatchesOnlyThisCheckoutsWatcher(t *testing.T) {
	pattern := regexp.MustCompile(watcherPattern("1143"))

	cases := []struct {
		name    string
		cmdline string
		want    bool
	}{
		{
			name:    "this slot's watcher",
			cmdline: "/home/dev/go/bin/gin --laddr localhost --port 3101 --appPort 1143 --build cmd/human --immediate --bin build/gin-bin -- --env-file .env serve",
			want:    true,
		},
		{
			name:    "the primary's watcher",
			cmdline: "/home/dev/go/bin/gin --laddr localhost --port 3001 --appPort 1043 --build cmd/human --immediate --bin build/gin-bin -- --env-file .env serve",
			want:    false,
		},
		{
			// 1143 is a prefix of 11430, and a slot's neighbour is not it.
			name:    "a port this one is a prefix of",
			cmdline: "/home/dev/go/bin/gin --laddr localhost --port 3101 --appPort 11430 --bin build/gin-bin",
			want:    false,
		},
		{
			name:    "the server child, which carries no appPort",
			cmdline: "/home/dev/Human/human/server/build/gin-bin --env-file .env serve",
			want:    false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := pattern.MatchString(c.cmdline); got != c.want {
				t.Fatalf("match(%q) = %v, want %v", c.cmdline, got, c.want)
			}
		})
	}
}

// TestDevServerPortIsTheURLsPort keeps the two readers of server/.env from
// drifting apart: the health check asks one port and the watcher lookup scopes
// itself by the other, and they have to be the same port.
func TestDevServerPortIsTheURLsPort(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "server"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "server", ".env"), []byte("HTTP_ADDR=:1543\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := devServerPort(root); got != "1543" {
		t.Fatalf("devServerPort() = %q, want 1543", got)
	}

	if got := devServerURL(root); got != "http://localhost:1543" {
		t.Fatalf("devServerURL() = %q, want http://localhost:1543", got)
	}
}
