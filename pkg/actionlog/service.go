package actionlog

import (
	"context"
	"strings"

	"github.com/go-chi/chi/middleware"

	"github.com/cortezaproject/corteza-server/pkg/api"
	"github.com/cortezaproject/corteza-server/pkg/auth"

	"go.uber.org/zap"
)

type (
	service struct {
		// where the audit log records are kept
		repo recordKeeper

		// Also write audit events here
		tee *zap.Logger

		// logger for repository errors
		logger *zap.Logger
	}

	loggable interface {
		LoggableAction() *Action
	}

	Recorder interface {
		Record(context.Context, loggable)
		Find(context.Context, Filter) (ActionSet, Filter, error)
	}

	recordKeeper interface {
		Record(context.Context, *Action) error
		Find(context.Context, Filter) (ActionSet, Filter, error)
	}
)

// NewService initializes auditlog service
//
func NewService(r recordKeeper, logger, tee *zap.Logger) (svc *service) {
	if tee == nil {
		tee = zap.NewNop()
	}

	svc = &service{
		tee:    tee,
		logger: logger,
		repo:   r,
	}

	return
}

func (svc service) Record(ctx context.Context, l loggable) {
	if l == nil {
		// nothing to record
		return
	}

	a := enrich(ctx, l.LoggableAction())

	var (
		log = svc.logger
	)

	zlf := []zap.Field{
		zap.Time("timestamp", a.Timestamp),
		zap.String("requestOrigin", a.RequestOrigin),
		zap.String("requestID", a.RequestID),
		zap.String("actorIPAddr", a.ActorIPAddr),
		zap.Uint64("actorID", a.ActorID),
		zap.String("resource", a.Resource),
		zap.String("action", a.Action),
		zap.Uint8("severity", uint8(a.Severity)),
		zap.String("error", a.Error),
		zap.String("description", a.Description),
		zap.Any("meta", a.Meta),
	}

	for k, v := range a.Meta {
		zlf = append(zlf, zap.Any("meta."+k, v))
	}

	log.Debug(a.Description, zlf...)

	if err := svc.repo.Record(ctx, a); err != nil {
		log.With(zap.Error(err)).Error("could not record audit event")
	}
}

func (svc service) Find(ctx context.Context, flt Filter) (ActionSet, Filter, error) {
	return svc.repo.Find(ctx, flt)
}

// Enriches action with additional info (ip, actor id, request id...)
func enrich(ctx context.Context, a *Action) *Action {
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
