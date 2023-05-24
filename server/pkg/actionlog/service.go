package actionlog

import (
	"context"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza/server/pkg/api"
	"github.com/cortezaproject/corteza/server/pkg/auth"
)

type (
	service struct {
		// where the audit log records are kept
		store actionlogStore

		// Also write audit events here
		tee *zap.Logger

		// logger for repository errors
		logger *zap.Logger

		policy policyMatcher
	}

	Recorder interface {
		Record(context.Context, *Action)
		Find(context.Context, Filter) (ActionSet, Filter, error)
	}

	actionlogStore interface {
		SearchActionlogs(ctx context.Context, f Filter) (ActionSet, Filter, error)
		CreateActionlog(ctx context.Context, rr ...*Action) error
	}
)

// NewService initializes action log service
func NewService(s actionlogStore, logger, tee *zap.Logger, policy policyMatcher) (svc *service) {
	if tee == nil {
		tee = zap.NewNop()
	}

	svc = &service{
		tee:    tee.Named("actionlog"),
		logger: logger.Named("actionlog"),
		store:  s,
		policy: policy,
	}

	return
}

func (svc service) Record(ctx context.Context, a *Action) {
	if a == nil {
		// nothing to record
		return
	}

	a = enrich(ctx, a)
	a.ID = id.Next()

	svc.log(a)
	if !svc.policy.Match(a) {
		// policy does not allow us to record this
		return
	}

	// We want to prevent any abrupt cancelation
	// (eg canceled request) that would cause
	// auditlog to fail...
	ctx = context.Background()

	if err := svc.store.CreateActionlog(ctx, a); err != nil {
		svc.logger.With(zap.Error(err)).Error("could not record audit event")
	}
}

func (svc service) log(a *Action) {
	zlf := []zap.Field{
		zap.Time("timestamp", a.Timestamp),
		zap.String("requestOrigin", a.RequestOrigin),
		zap.String("requestID", a.RequestID),
		zap.String("actorIPAddr", a.ActorIPAddr),
		logger.Uint64("actorID", a.ActorID),
		zap.String("resource", a.Resource),
		zap.String("action", a.Action),
		zap.Uint8("severity", uint8(a.Severity)),
		zap.String("error", a.Error),
		zap.String("description", a.Description),
		zap.Bool("policy-match", svc.policy.Match(a)),
		zap.Any("meta", a.Meta),
	}

	svc.tee.
		With(zlf...).
		// Skipping 3 callers (the most common stack)
		//   actionlog.service.log()
		//   actionlog.service.Record()
		//   (generated service function)
		//
		// One exception, access control, that calls Record fn directly,
		// without going through generated actionlog helpers
		WithOptions(zap.AddCallerSkip(3)).
		// This is debug logger and we log all recordings as debug
		Debug(a.Description)
}

func (svc service) Find(ctx context.Context, flt Filter) (ActionSet, Filter, error) {
	return svc.store.SearchActionlogs(ctx, flt)
}

// Enriches action with additional info (ip, actor id, request id...)
func enrich(ctx context.Context, a *Action) *Action {
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}

	a.RequestOrigin = RequestOriginFromContext(ctx)

	// Relies on chi's middleware to get to the request ID
	// This does not hurt us for now.
	a.RequestID = middleware.GetReqID(ctx)

	// uses pkg/auth to extract stored identity from context
	a.ActorID = auth.GetIdentityFromContext(ctx).Identity()

	// IP from the request,
	// we're splitting by space & colon to remove any additional (proxy) IPs
	// and ports from the string
	if tmp := strings.SplitN(api.RemoteAddrFromContext(ctx), " ", 2); len(tmp) > 0 {
		// split by : (ip:port)
		if tmp = strings.SplitN(tmp[0], ":", 2); len(tmp) > 0 {
			const maxLen = 16
			ipAddr := tmp[0]
			if len(ipAddr) > maxLen {
				ipAddr = ipAddr[:maxLen-1]
			}

			a.ActorIPAddr = ipAddr
		}
	}

	return a
}
