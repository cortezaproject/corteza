package provision

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/service"
	"go.uber.org/zap"
)

func Run(ctx context.Context, log *zap.Logger, s store.Storer, path string) error {
	ffn := []func() error{
		func() error { return roles(ctx, s) },
		func() error { return importConfig(ctx, log, s, path) },
		func() error { return authSettingsAutoDiscovery(ctx, log, service.DefaultSettings) },
		func() error { return authAddExternals(ctx, log) },
		func() error { return service.DefaultSettings.UpdateCurrent(ctx) },
		func() error { return oidcAutoDiscovery(ctx, log) },
	}

	for _, fn := range ffn {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}
