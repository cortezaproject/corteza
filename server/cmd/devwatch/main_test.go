package main

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

// TestRebuildServesAChangeRaisedDuringABuild is the whole reason this command
// exists, and the one behaviour a replacement must not lose.
//
// The watcher this replaces stamped a clock after each build returned, so a
// file written during those ~15s was already older than the stamp and was never
// picked up — an edit silently skipped for good, which reads as a fix that did
// not work. Here the write becomes a pending signal instead, and a pending
// signal survives the build it arrived during.
func TestRebuildServesAChangeRaisedDuringABuild(t *testing.T) {
	changed := make(chan struct{}, 1)
	stop := make(chan struct{})

	var (
		mu      sync.Mutex
		builds  int
		stopped bool
	)

	// Bounded, because the failure this guards against is not a wrong count but
	// a loop with nothing left to wake it: a dropped signal means the second
	// build never comes, and an unbounded rebuild would hang rather than fail.
	runRebuild(t, changed, stop, func() {
		mu.Lock()
		builds++
		n := builds
		mu.Unlock()

		switch n {
		case 1:
			// A write landing mid-build, which is the case that was lost.
			raise(changed)
		case 2:
			mu.Lock()
			stopped = true
			mu.Unlock()
			close(stop)
		}
	})

	mu.Lock()
	defer mu.Unlock()

	if builds != 2 {
		t.Fatalf("builds = %d, want 2 — the change raised during the first build was dropped", builds)
	}

	if !stopped {
		t.Fatal("the loop ended without reaching the second build")
	}
}

// TestRebuildCoalescesAStormOfWrites keeps a multi-file save from costing a
// build per file. A single rebuild that includes every write is the correct
// answer, and it is also the fast one.
func TestRebuildCoalescesAStormOfWrites(t *testing.T) {
	changed := make(chan struct{}, 1)
	stop := make(chan struct{})

	var (
		mu     sync.Mutex
		builds int
	)

	runRebuild(t, changed, stop, func() {
		mu.Lock()
		builds++
		n := builds
		mu.Unlock()

		if n == 1 {
			for range 20 {
				raise(changed)
			}
			return
		}

		close(stop)
	})

	mu.Lock()
	defer mu.Unlock()

	if builds != 2 {
		t.Fatalf("builds = %d, want 2 — twenty writes must not be twenty builds", builds)
	}
}

// TestRebuildStopsWithoutAnyChange pins that stopping does not need a write to
// notice it: `make watch` interrupted on an idle tree has to exit.
func TestRebuildStopsWithoutAnyChange(t *testing.T) {
	changed := make(chan struct{}, 1)
	stop := make(chan struct{})
	close(stop)

	done := make(chan struct{})

	go func() {
		rebuild(changed, stop, 0, func() {})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("rebuild did not return on a closed stop channel")
	}
}

// TestRaiseNeverBlocks guards the property that lets the event goroutine record
// a change without waiting for the build to finish. A blocking send there would
// stall the watcher, and fsnotify drops events it cannot deliver.
func TestRaiseNeverBlocks(t *testing.T) {
	changed := make(chan struct{}, 1)

	done := make(chan struct{})

	go func() {
		for range 100 {
			raise(changed)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("raise blocked on a full channel")
	}

	if len(changed) != 1 {
		t.Fatalf("pending signals = %d, want 1", len(changed))
	}
}

// TestWatchTreeWatchesTheRootWhateverItIsCalled pins the case that makes a
// watcher watch nothing at all.
//
// The skip list rejects dotted directories, and the walk visits the root first.
// Applied to the root, a relative root of "." is skipped along with everything
// under it — and a watcher watching zero directories behaves exactly like one
// with nothing to do.
func TestWatchTreeWatchesTheRootWhateverItIsCalled(t *testing.T) {
	root := t.TempDir()

	for _, dir := range []string{"compose/service", "pkg/dal", "vendor/example.com/lib", "build", "node_modules/x", ".git/objects"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	watcher := newTestWatcher(t)

	added, err := watchTree(watcher, root)
	if err != nil {
		t.Fatalf("watchTree: %v", err)
	}

	// root, compose, compose/service, pkg, pkg/dal — and nothing excluded.
	if added != 5 {
		t.Fatalf("watched %d directories, want 5 (root + compose + compose/service + pkg + pkg/dal)", added)
	}

	// The same tree reached by a relative "." must come out the same, which is
	// what the root exemption buys.
	t.Chdir(root)

	relative, err := watchTree(newTestWatcher(t), ".")
	if err != nil {
		t.Fatalf("watchTree(.): %v", err)
	}

	if relative != added {
		t.Fatalf("watching \".\" found %d directories, absolute found %d", relative, added)
	}
}

func TestSkipDir(t *testing.T) {
	skipped := []string{"vendor", "build", "node_modules", ".git", "testdata", "var", ".idea"}
	kept := []string{"compose", "pkg", "cmd", "system", "automation"}

	for _, name := range skipped {
		if !skipDir(name) {
			t.Errorf("skipDir(%q) = false, want true", name)
		}
	}

	for _, name := range kept {
		if skipDir(name) {
			t.Errorf("skipDir(%q) = true, want false", name)
		}
	}
}

// TestIsSource keeps an editor's own scratch files from costing a rebuild each.
func TestIsSource(t *testing.T) {
	sources := []string{"record.go", "/abs/path/record_test.go"}
	noise := []string{
		"record.go~",              // emacs backup
		"/abs/path/.record.go.sw", // vim swap
		"record.md",
		"schema.yaml",
		"/abs/path/build/dev-bin",
	}

	for _, name := range sources {
		if !isSource(name) {
			t.Errorf("isSource(%q) = false, want true", name)
		}
	}

	for _, name := range noise {
		if isSource(name) {
			t.Errorf("isSource(%q) = true, want false", name)
		}
	}
}

// runRebuild runs the loop to completion, or fails the test rather than hanging
// on one that has nothing left to wake it.
func runRebuild(t *testing.T, changed <-chan struct{}, stop <-chan struct{}, cycle func()) {
	t.Helper()

	done := make(chan struct{})

	go func() {
		rebuild(changed, stop, 0, cycle)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("rebuild never finished — a raised change was dropped, so nothing woke the loop")
	}
}

func newTestWatcher(t *testing.T) *fsnotify.Watcher {
	t.Helper()

	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("cannot create a watcher: %v", err)
	}

	t.Cleanup(func() { _ = w.Close() })

	return w
}
