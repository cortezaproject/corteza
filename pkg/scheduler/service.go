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
		active chan bool
	}

	dispatcher interface {
		Dispatch(ctx context.Context, ev eventbus.Event)
	}
)

const (
	defaultInterval = time.Minute
)

var (
	now = func() time.Time { return time.Now() }

	// Global scheduler
	gScheduler *service
)

// Setup configures global scheduling service
func Setup(log *zap.Logger, d dispatcher, interval time.Duration) {
	if gScheduler != nil {
		// shut it down
		gScheduler.active <- false
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
	}

	return svc
}

// Register all events that should fire on tick (interval)
func (svc *service) OnTick(events ...eventbus.Event) {
	svc.events = append(svc.events, events...)
}

// Run starts event scheduler service
func (svc service) Start(ctx context.Context) {
	if svc.active != nil {
		return
	}

	svc.active = make(chan bool, 1)
	go func() {
		defer sentry.Recover()

		nextTick := now().Truncate(svc.interval).Add(svc.interval)

		svc.log.Info(
			"starting",
			zap.Time("delay", nextTick),
			zap.Duration("interval", svc.interval),
		)

		// Wait until start of the next interval
		time.Sleep(nextTick.Sub(now()))
		svc.log.Info("started")

		// start with first interval
		svc.dispatch(ctx)
		ticker := time.NewTicker(svc.interval)
		defer ticker.Stop()
		defer svc.log.Info("stopped")

		for {
			select {
			case <-svc.active:
				svc.log.Info("unactivated")
				return
			case <-ctx.Done():
				svc.log.Info("done")
				return
			case <-ticker.C:
				svc.dispatch(ctx)
			}
		}

	}()
}

func (svc service) dispatch(ctx context.Context) {
	for _, ev := range svc.events {
		go func(ev eventbus.Event) {
			sentry.Recover()
			svc.dispatcher.Dispatch(ctx, ev)
		}(ev)
	}
}
