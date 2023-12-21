package provision

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/store"
	"go.uber.org/zap"
)

var (
	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}
)

func Run(ctx context.Context, log *zap.Logger, s store.Storer, provisionOpt options.ProvisionOpt, authOpt options.AuthOpt) error {
	log = log.Named("provision")

	// Note,
	ffn := []func() error{
		// Migrations:
		// (placeholder for all post 2022.9.x modifications)
		func() error { return migratePages(ctx, log.Named("pages"), s) },

		// *************************************************************************************************************

		// Config (full & partial)
		func() error { return importConfig(ctx, log.Named("config"), s, provisionOpt.Path) },

		// *************************************************************************************************************

		// Auto-discoveries and other parts that cannot be imported from static files
		func() error { return emailSettings(ctx, s) },
		func() error { return apigwFilters(ctx, log.Named("apigw.filters"), s) },
		func() error { return authAddExternals(ctx, log.Named("auth.externals"), s) },
		func() error { return oidcAutoDiscovery(ctx, log.Named("auth.oidc-auto-discovery"), s, authOpt) },
		func() error { return defaultAuthClient(ctx, log.Named("auth.clients"), s, authOpt) },
		func() error { return addAuthSuperUsers(ctx, log.Named("auth.super-users"), s, authOpt) },
		func() error { return invalidateDedupRules(ctx, log.Named("compose.deduplication"), s) },
	}

	for _, fn := range ffn {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
