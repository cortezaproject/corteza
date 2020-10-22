package store

import (
	"context"
	"go.uber.org/zap"
)

type (
	storeUpgrader interface {
		Upgrade(context.Context, *zap.Logger) error
	}
)

func Upgrade(ctx context.Context, log *zap.Logger, s Storer) error {
	upgradableStore, ok := s.(storeUpgrader)
	if !ok {
		log.Debug("store does not support upgrades")
		return nil
	}

	return upgradableStore.Upgrade(ctx, log)
}
