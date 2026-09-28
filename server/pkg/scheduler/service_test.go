package scheduler

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza/server/pkg/eventbus"
)

type (
	mockEvent struct {
		rType string
		eType string
		match func(matcher eventbus.ConstraintMatcher) bool
	}
)

func (e mockEvent) ResourceType() string {
	return e.rType
}

func (e mockEvent) EventType() string {
	return e.eType
}

func (e mockEvent) Match(matcher eventbus.ConstraintMatcher) bool {
	if e.match == nil {
		return true
	}

	return e.match(matcher)
}

func TestMainServiceFunctions(t *testing.T) {
	r := require.New(t)

	const (
		loopInterval = time.Millisecond * 100
		actionWait   = loopInterval * 10
	)

	if gScheduler != nil {
		gScheduler.Stop()
		gScheduler = nil
	}

	Setup(zap.NewNop(), eventbus.New(), loopInterval)
	r.NotNil(gScheduler)
	r.False(gScheduler.Started())
	r.Equal(gScheduler, Service())
	Service().Start(context.Background())
	Service().OnTick(&mockEvent{}, &mockEvent{}, &mockEvent{}, &mockEvent{})
	time.Sleep(actionWait)
	r.True(gScheduler.Started())
	gScheduler.Stop()
	time.Sleep(actionWait)
	r.False(gScheduler.Started())
}

func TestDispatchRunsHandlersIndependently(t *testing.T) {
	var (
		r   = require.New(t)
		bus = eventbus.New()
		svc = NewService(zap.NewNop(), bus, time.Minute)
		ev  = &mockEvent{rType: "system", eType: "onInterval"}

		wg      sync.WaitGroup
		started = make(chan int, 3)
		release = make(chan struct{})
	)

	wg.Add(3)

	bus.Register(func(context.Context, eventbus.Event) error {
		defer wg.Done()
		started <- 0
		return fmt.Errorf("failed")
	}, eventbus.For("system"), eventbus.On("onInterval"), eventbus.Weight(0))

	for w := 1; w <= 2; w++ {
		w := w
		bus.Register(func(context.Context, eventbus.Event) error {
			defer wg.Done()
			started <- w
			<-release
			return nil
		}, eventbus.For("system"), eventbus.On("onInterval"), eventbus.Weight(w))
	}

	svc.OnTick(ev)
	svc.dispatch(context.Background())

	seen := map[int]bool{}
	for len(seen) < 3 {
		select {
		case w := <-started:
			seen[w] = true
		case <-time.After(time.Second):
			r.FailNow("handlers did not all start", "started: %v", seen)
		}
	}

	close(release)
	wg.Wait()
}
