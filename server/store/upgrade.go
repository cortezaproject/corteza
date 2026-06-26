package store

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Upgrade runs all needed upgrades on a specific store
func Upgrade(ctx context.Context, log *zap.Logger, s Storer) error {
	s.SetLogger(log)
	return s.Upgrade(ctx)
}

type actionlogUpgrader interface {
	UpgradeActionlog(ctx context.Context) error
}

// UpgradeActionlog creates only the tables required by the actionlog service.
func UpgradeActionlog(ctx context.Context, log *zap.Logger, s Storer) error {
	u, ok := s.(actionlogUpgrader)
	if !ok {
		return fmt.Errorf("store does not support actionlog upgrade")
	}
	s.SetLogger(log)
	return u.UpgradeActionlog(ctx)
}
