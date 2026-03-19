package governor

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"go.uber.org/zap"
)

var testIDCounter atomic.Uint64

func nextID() id.ID {
	return id.MustNumID(testIDCounter.Add(1))
}

// newTestGovernor returns a governor with a fake clock so tests can control time.
func newTestGovernor(t *testing.T) (*governor, *time.Time) {
	t.Helper()
	now := time.Now()
	g := &governor{
		gates: gates{
			globalPause: newGate(),
			exec:        make(map[id.ID]*execGates),
		},
		exec: make(map[id.ID]*execPolicy),
		config: Config{
			WatcherFallbackInterval: 5 * time.Second,
		},
		now: func() time.Time { return now },
	}
	g.log = zap.NewNop()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go g.watch(ctx)

	return g, &now
}

func isClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}
