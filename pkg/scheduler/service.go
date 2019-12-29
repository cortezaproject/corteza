package scheduler

import (
	"context"
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

		// Simple chan to control if service is running or not
		ticker *time.Ticker
	}

	dispatcher interface {
		Dispatch(ctx context.Context, ev eventbus.Event)
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
		log:        log.Named("scheduler"),
		interval:   interval,
		dispatcher: d,
		events:     make([]eventbus.Event, 0, maxEvents),
	}

	return svc
}

// Register all events that should fire on tick (interval)
func (svc *service) OnTick(events ...eventbus.Event) {
	svc.events = append(svc.events, events...)
}

func (svc *service) Stop() {
	if svc.ticker == nil {
		svc.log.Debug("already stopped")
	} else {
		svc.log.Debug("stopping")
		svc.ticker.Stop()
		svc.ticker = nil
	}
}

// Run starts event scheduler service
func (svc *service) Start(ctx context.Context) {
	if svc.ticker != nil {
		svc.log.Debug("already started")
		return
	}

	svc.ticker = &time.Ticker{}

	go func() {
		defer sentry.Recover()

		nextTick := now().Truncate(svc.interval).Add(svc.interval)

		svc.log.Debug(
			"starting",
			zap.Time("delay", nextTick),
			zap.Duration("interval", svc.interval),
		)

		// Wait until start of the next interval
		time.Sleep(nextTick.Sub(now()))
		svc.log.Debug("started")

		// start with first interval
		svc.dispatch(ctx)
		svc.ticker = time.NewTicker(svc.interval)

	loop:
		for {
			select {
			case <-svc.ticker.C:
				svc.dispatch(ctx)

			case <-ctx.Done():
				svc.log.Debug("done")
				break loop
			}
		}

		defer svc.log.Debug("stopped")
		svc.ticker.Stop()
		svc.ticker = nil
	}()
}

func (svc service) Started() (started bool) {
	return svc.ticker != nil
}

func (svc service) dispatch(ctx context.Context) {
	for _, ev := range svc.events {
		go func(ev eventbus.Event) {
			sentry.Recover()
			svc.dispatcher.Dispatch(ctx, ev)
		}(ev)
	}
}
