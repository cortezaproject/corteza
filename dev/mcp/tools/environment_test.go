package tools

import (
	"os"
	"path/filepath"
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
			// The regression itself: mtimes say fresh, the process says otherwise.
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
