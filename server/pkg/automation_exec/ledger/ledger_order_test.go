package ledger

import (
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/automation_exec/types"
	"github.com/crusttech/human/server/pkg/id"
)

// stamp overwrites an execution's CreatedAt so ordering can be asserted without
// sleeping between registrations.
func stamp(t *testing.T, l *ledger, xID, eID id.ID, rev int, at time.Time) {
	t.Helper()
	ex, ok := l.store[xID][eID][rev]
	if !ok {
		t.Fatalf("execution %v not registered", eID)
	}
	ex.CreatedAt = at
}

// TL;DR: both list calls return executions newest first.
// Example: an agent fires a probe record, then reads executions[0] to find the
// run it just caused. The store is a map, so unordered iteration made that the
// latest run only by chance — and sent whoever read a trace from it to an
// arbitrary execution.
func TestListExecutions_NewestFirst(t *testing.T) {
	base := time.Date(2026, 8, 25, 20, 0, 0, 0, time.UTC)

	// Registered oldest-to-newest, so a stable-but-unsorted implementation
	// would tend to return exactly the reverse of what is expected.
	newLedgerWith := func(t *testing.T) (*ledger, id.ID, []id.ID) {
		t.Helper()
		l := newLedger()
		xID := nextID()
		var eIDs []id.ID
		for i := 0; i < 8; i++ {
			eID := nextID()
			if err := l.RegisterExecution(ctx, xID, eID, 1, types.ExecutionParams{}); err != nil {
				t.Fatal(err)
			}
			stamp(t, l, xID, eID, 1, base.Add(time.Duration(i)*time.Minute))
			eIDs = append(eIDs, eID)
		}
		return l, xID, eIDs
	}

	assertDescending := func(t *testing.T, got []*types.Execution, want []id.ID) {
		t.Helper()
		if len(got) != len(want) {
			t.Fatalf("expected %d executions, got %d", len(want), len(got))
		}
		for i := range want {
			// want is oldest-first; the result must be its reverse.
			expect := want[len(want)-1-i]
			if got[i].ID != expect {
				t.Fatalf("position %d: expected %v, got %v", i, expect, got[i].ID)
			}
		}
	}

	t.Run("ListExecutionsByExecutable", func(t *testing.T) {
		l, xID, eIDs := newLedgerWith(t)
		got, err := l.ListExecutionsByExecutable(ctx, xID)
		if err != nil {
			t.Fatal(err)
		}
		assertDescending(t, got, eIDs)
	})

	t.Run("ListExecutions", func(t *testing.T) {
		l, _, eIDs := newLedgerWith(t)
		got, err := l.ListExecutions(ctx)
		if err != nil {
			t.Fatal(err)
		}
		assertDescending(t, got, eIDs)
	})

	// Repeated calls must agree. A single call can match by luck; two calls
	// disagreeing is what the map actually did.
	t.Run("stable across calls", func(t *testing.T) {
		l, xID, _ := newLedgerWith(t)
		first, _ := l.ListExecutionsByExecutable(ctx, xID)
		for n := 0; n < 20; n++ {
			again, _ := l.ListExecutionsByExecutable(ctx, xID)
			for i := range first {
				if first[i].ID != again[i].ID {
					t.Fatalf("call %d differs at position %d: %v then %v", n, i, first[i].ID, again[i].ID)
				}
			}
		}
	})
}

// TL;DR: executions sharing a timestamp fall back to ID descending.
// Example: two triggers on one record write are registered in the same
// instant; without a tie-break their relative order is still the map's.
func TestListExecutions_TieBreaksOnID(t *testing.T) {
	l := newLedger()
	xID := nextID()
	at := time.Date(2026, 8, 25, 20, 0, 0, 0, time.UTC)

	var eIDs []id.ID
	for i := 0; i < 8; i++ {
		eID := nextID()
		if err := l.RegisterExecution(ctx, xID, eID, 1, types.ExecutionParams{}); err != nil {
			t.Fatal(err)
		}
		stamp(t, l, xID, eID, 1, at)
		eIDs = append(eIDs, eID)
	}

	got, err := l.ListExecutionsByExecutable(ctx, xID)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(eIDs) {
		t.Fatalf("expected %d, got %d", len(eIDs), len(got))
	}
	for i := 1; i < len(got); i++ {
		if got[i-1].ID.Num() <= got[i].ID.Num() {
			t.Fatalf("position %d: IDs not descending (%v then %v)", i, got[i-1].ID, got[i].ID)
		}
	}
}
