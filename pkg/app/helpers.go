package app

import (
	"context"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
)

// RunSetup calls Setup hooks on all runnable parts
//
// It stops on first error
func RunSetup(log *zap.Logger, opts *Options, pp ...Runnable) (err error) {
	for _, app := range pp {
		err = app.Setup(log, opts)
		if err != nil {
			return
		}
	}

	return
}

// RunInitialize calls Initialize hooks on all runnable parts
//
// It stops on first error
func RunInitialize(ctx context.Context, pp ...Runnable) (err error) {
	for _, app := range pp {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Init)
		err = app.Initialize(ctx)
		if err != nil {
			return
		}
	}

	return
}

// RunUpgrade calls Upgrade hooks on all runnable parts
//
// It stops on first error
func RunUpgrade(ctx context.Context, pp ...Runnable) (err error) {
	for _, app := range pp {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Upgrade)
		err = app.Upgrade(ctx)
		if err != nil {
			return
		}
	}

	return
}

// RunActivate calls Activate hooks on all runnable parts
//
// It stops on first error
func RunActivate(ctx context.Context, pp ...Runnable) (err error) {
	for _, app := range pp {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Activate)
		err = app.Activate(ctx)
		if err != nil {
			return
		}
	}

	return
}

// RunProvision calls Provision hooks on all runnable parts
//
// It stops on first error
func RunProvision(ctx context.Context, pp ...Runnable) (err error) {
	for _, app := range pp {
		ctx = actionlog.RequestOriginToContext(ctx, actionlog.RequestOrigin_APP_Provision)
		err = app.Provision(ctx)
		if err != nil {
			return
		}
	}

	return
}
