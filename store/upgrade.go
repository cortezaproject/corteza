package store

import (
	"context"

	"go.uber.org/zap"
)

// Upgrade runs all needed upgrades on a specific store
func Upgrade(ctx context.Context, log *zap.Logger, s Storer) error {
	return Tx(ctx, s, func(ctx context.Context, s Storer) error {
		s.SetLogger(log)
		return s.Upgrade(ctx)
	})
}
