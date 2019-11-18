package service

import (
	"context"

	"github.com/titpetric/factory"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	statistics struct {
		db     *factory.DB
		ctx    context.Context
		logger *zap.Logger

		ac statisticsAccessController
	}

	statisticsAccessController interface {
		CanAccess(context.Context) bool
	}
)

func Statistics(ctx context.Context) *statistics {
	return &statistics{
		db:     repository.DB(ctx),
		ac:     DefaultAccessControl,
		logger: DefaultLogger.Named("statistics"),
	}
}

// log() returns zap's logger with requestID from current context and fields.
func (svc statistics) log(ctx context.Context, fields ...zapcore.Field) *zap.Logger {
	return logger.AddRequestID(ctx, svc.logger).With(fields...)
}

func (svc statistics) Metrics(ctx context.Context) (interface{}, error) {
	if !svc.ac.CanAccess(ctx) {
		return nil, ErrNoPermissions
	}

	type (
		metricsPayload struct {
			Users        *types.UserMetrics        `json:"users"`
			Roles        *types.RoleMetrics        `json:"roles"`
			Applications *types.ApplicationMetrics `json:"applications"`
		}
	)

	var (
		rval = &metricsPayload{}
		err  error
	)

	if rval.Users, err = repository.User(ctx, svc.db).Metrics(); err != nil {
		return nil, err
	}

	if rval.Roles, err = repository.Role(ctx, svc.db).Metrics(); err != nil {
		return nil, err
	}

	if rval.Applications, err = repository.Application(ctx, svc.db).Metrics(); err != nil {
		return nil, err
	}

	return rval, err
}
