package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/crusttech/human/server/pkg/eventbus"
)

// recordingDispatcher notes when each dispatch happened
type recordingDispatcher struct {
	mux   sync.Mutex
	times []time.Time
}

func (d *recordingDispatcher) WaitForAll(_ context.Context, _ eventbus.Event) []error {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.times = append(d.times, time.Now())
	return nil
}

func (d *recordingDispatcher) snapshot() []time.Time {
	d.mux.Lock()
	defer d.mux.Unlock()
	return append([]time.Time(nil), d.times...)
}

// Every cycle must land on an interval boundary of the wall clock; a ticker
// measured from the previous delivery let the dispatch creep away from the
// boundary until scheduled automations stopped matching their interval.
func TestWatchAlignsToIntervalBoundary(t *testing.T) {
	const (
		interval  = 200 * time.Millisecond
		tolerance = 60 * time.Millisecond
	)

	r := require.New(t)
	d := &recordingDispatcher{}
	svc := NewService(zap.NewNop(), d, interval)
	svc.OnTick(&mockEvent{})

	ctx, cancel := context.WithCancel(context.Background())
	svc.Start(ctx)
	time.Sleep(interval*5 + interval/2)
	cancel()
	svc.Stop()

	// the first dispatch happens right at start and is not aligned
	times := d.snapshot()
	r.GreaterOrEqual(len(times), 5)

	for _, ts := range times[1:] {
		offset := ts.Sub(ts.Truncate(interval))
		r.Less(offset, tolerance, "dispatch at %s is %s past the interval boundary", ts.Format("15:04:05.000"), offset)
	}
}
