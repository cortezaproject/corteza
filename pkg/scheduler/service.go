package scheduler

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/pkg/sentry"
)

type (
	service struct {
		log        *zap.Logger
		events     []eventbus.Event
		interval   time.Duration
		dispatcher dispatcher

		// Read & write locking
		l sync.RWMutex

		// Simple chan to control if service is running or not
		ticker *time.Ticker
	}

	dispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) error
	}
)

const (
	defaultInterval = time.Minute

	// There should not be more than 2 per each service (<no of services> * 2 [interval + timestamp])
	maxEvents = 16
)

var (
	now = func() time.Time { return time.Now() }

	// Global scheduler
	gScheduler *service
)

// Setup configures global scheduling service
func Setup(log *zap.Logger, d dispatcher, interval time.Duration) {
	if gScheduler != nil {
		// shut down current global scheduler
		gScheduler.Stop()
	}

	gScheduler = NewService(log, d, interval)
}

func Service() *service {
	return gScheduler
}

func NewService(log *zap.Logger, d dispatcher, interval time.Duration) *service {
	// Fix interval to positive number
	if interval == 0 {
		interval = defaultInterval
	}

	var svc = &service{
		l:          sync.RWMutex{},
		log:        log.Named("scheduler"),
		interval:   interval,
		dispatcher: d,
		events:     make([]eventbus.Event, 0, maxEvents),
	}

	return svc
}

// Register all events that should fire on tick (interval)
func (svc *service) OnTick(events ...eventbus.Event) {
	svc.l.Lock()
	defer svc.l.Unlock()
	svc.events = append(svc.events, events...)
}

func (svc *service) Stop() {

	if svc.ticker == nil {
		svc.log.Debug("already stopped")
	} else {
		svc.log.Debug("stopping")
		svc.ticker.Stop()
		svc.l.Lock()
		svc.ticker = nil
		defer svc.l.Unlock()
	}
}

// Run starts event scheduler service
func (svc *service) Start(ctx context.Context) {

	if svc.ticker != nil {
		svc.log.Debug("already started")
		return
	}

	svc.l.Lock()
	svc.ticker = &time.Ticker{}
	svc.l.Unlock()

	go func() {
		defer sentry.Recover()

		// Calculate how much time we need to wait until next tick
		delay := now().Truncate(svc.interval).Add(svc.interval).Sub(now())

		svc.log.Debug(
			"starting",
			zap.Duration("delay", delay),
			zap.Duration("interval", svc.interval),
		)

		// Wait until start of the next interval
		time.Sleep(delay)
		svc.l.Lock()
		svc.ticker = time.NewTicker(svc.interval)
		svc.l.Unlock()
		svc.log.Debug("started")

		go svc.watch(ctx)
	}()
}

func (svc *service) watch(ctx context.Context) {
	defer sentry.Recover()
	defer func() {
		defer svc.log.Debug("stopped")
		svc.ticker.Stop()
		svc.l.Lock()
		svc.ticker = nil
		svc.l.Unlock()
	}()

	// start with first interval
	svc.dispatch(ctx)

	for {
		select {
		case <-svc.ticker.C:
			svc.dispatch(ctx)

		case <-ctx.Done():
			svc.log.Debug("done")
			return
		}
	}
}

func (svc *service) Started() (started bool) {
	svc.l.RLock()
	defer svc.l.RUnlock()

	return svc.ticker != nil
}

func (svc *service) dispatch(ctx context.Context) {
	svc.l.RLock()

	ee := make([]eventbus.Event, len(svc.events))
	for e := range svc.events {
		ee[e] = svc.events[e]
	}

	defer svc.l.RUnlock()

	for _, ev := range ee {
		go func(ev eventbus.Event) {
			err := svc.dispatcher.WaitFor(ctx, ev)
			if err != nil {
				svc.log.Warn("failed to execute scheduled trigger", zap.Error(err))
			}
		}(ev)
	}
}
