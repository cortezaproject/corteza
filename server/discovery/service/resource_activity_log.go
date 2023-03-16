package service

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/eventbus"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"go.uber.org/zap"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/discovery/types"
	"github.com/cortezaproject/corteza/server/store"
)

type (
	resourceActivity struct {
		// log for repository errors
		log *zap.Logger
		opt options.DiscoveryOpt
		mux sync.RWMutex

		// where the activity log records are kept
		// store resourceActivityLogStore
		store store.Storer

		eventbus eventbusRegistry
		decoder  types.ResDecoder
	}

	eventbusRegistry interface {
		Register(eventbus.HandlerFn, ...eventbus.HandlerRegOp) uintptr
	}

	ResourceActivityService interface {
		Find(context.Context, types.ResourceActivityFilter) (types.ResourceActivitySet, types.ResourceActivityFilter, error)
	}
)

func ResourceActivity(log *zap.Logger, opt options.DiscoveryOpt, s store.Storer, eb eventbusRegistry) *resourceActivity {
	return &resourceActivity{
		log:      log.Named("resource-activity"),
		opt:      opt,
		store:    s,
		eventbus: eb,
	}
}

func (svc *resourceActivity) logger(a *types.ResourceActivity) {
	zlf := []zap.Field{
		zap.Uint8("recordID", uint8(a.ResourceID)),
		zap.String("ResourceType", a.ResourceType),
		zap.String("ResourceAction", a.ResourceAction),
		zap.Time("timestamp", a.Timestamp),
	}

	svc.log.
		With(zlf...).
		// Skipping 3 callers (the most common stack)
		//   discovery.service.log()
		//   discovery.service.Record()
		//   (generated service function)
		//
		// One exception, access control, that calls Record fn directly,
		// without going through generated discovery helpers
		WithOptions(zap.AddCallerSkip(3)).
		// This is debug logger, and we log all recordings as debug
		Debug(fmt.Sprintf("%s of %s", a.ResourceAction, a.ResourceType))
}

func (svc *resourceActivity) InitResourceActivityLog(ctx context.Context, resourceType []string) (err error) {
	eventType := []string{
		string(types.AfterCreate),
		string(types.AfterUpdate),
		string(types.AfterDelete),
	}

	svc.mux.RLock()
	defer svc.mux.RUnlock()

	svc.eventbus.Register(
		func(_ context.Context, ev eventbus.Event) error {
			var a *types.ResourceActivity
			dec, is := ev.(types.ResDecoder)
			if is {
				svc.log.Debug("resource changed",
					zap.String("eventType", ev.EventType()),
					zap.String("resourceType", ev.ResourceType()),
				)

				a, err = types.CastToResourceActivity(dec)
				if err != nil {
					svc.log.With(zap.Error(err)).Error("could not cast event to activity")
					return err
				}

				err = svc.store.CreateResourceActivity(ctx, a)
				if err != nil {
					svc.log.With(zap.Error(err)).Error("could not record activity event")
					return err
				}

				return nil
			}

			return nil
		},
		eventbus.For(resourceType...),
		eventbus.On(eventType...),
	)

	return
}

func (svc *resourceActivity) Find(ctx context.Context, filter types.ResourceActivityFilter) (aa types.ResourceActivitySet, f types.ResourceActivityFilter, err error) {
	err = func() error {
		aa, f, err = store.SearchResourceActivitys(ctx, svc.store, filter)
		if err != nil {
			return err
		}

		return nil
	}()

	return aa, f, err
}

func (svc *resourceActivity) Record(ctx context.Context, a *types.ResourceActivity) {
	if a == nil {
		// nothing to record
		return
	}

	a = enrich(ctx, a)
	a.ID = id.Next()

	svc.logger(a)

	// We want to prevent any abrupt cancellation
	// (e.g. canceled request) that would cause
	// discovery to fail...
	ctx = context.Background()

	if err := svc.store.CreateResourceActivity(ctx, a); err != nil {
		svc.log.With(zap.Error(err)).Error("could not record activity event")
	}
}

// enrich activity with additional info (timestamp, ...)
func enrich(_ context.Context, a *types.ResourceActivity) *types.ResourceActivity {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	return a
}
