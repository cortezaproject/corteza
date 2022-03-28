package store

import (
	"context"

	"go.uber.org/zap"
)

type (
	upgradable interface {
		Upgrade(context.Context, *zap.Logger) error
	}
)

// Upgrade runs all needed upgrades on a specific store
//
// See store/adapters/rdbms/schema/README.md for more details on RDBMS implementations
func Upgrade(ctx context.Context, log *zap.Logger, s upgradable) error {
	return s.Upgrade(ctx, log)
}
