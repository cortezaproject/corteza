package rdbms

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/schema"
	"go.uber.org/zap"
)

func (s Store) Upgrade(ctx context.Context, log *zap.Logger) error {
	if s.config == nil {
		return fmt.Errorf("config not set on RDBMS store")
	}

	if s.config.Upgrader == nil {
		return fmt.Errorf("upgrader not configured on RDBMS store")

	}

	s.config.Upgrader.SetLogger(log)

	return schema.Upgrade(ctx, s.config.Upgrader)
}
