package cred_registry

import (
	"context"
	"time"

	"go.uber.org/zap"
)

const (
	// RefreshCheckInterval is how often the refresher checks for expiring tokens
	RefreshCheckInterval = 1 * time.Minute
)

type (
	refresher struct {
		registry *registry
		ticker   *time.Ticker
		done     chan bool
		logger   *zap.Logger
	}
)

func newRefresher(reg *registry, logger *zap.Logger) *refresher {
	return &refresher{
		registry: reg,
		ticker:   time.NewTicker(RefreshCheckInterval),
		done:     make(chan bool),
		logger:   logger,
	}
}

func (rf *refresher) start() {
	go rf.run()
}

func (rf *refresher) stop() {
	rf.ticker.Stop()
	rf.done <- true
}

func (rf *refresher) run() {
	rf.logger.Info("credential refresher started")

	for {
		select {
		case <-rf.ticker.C:
			rf.refreshAll()
		case <-rf.done:
			rf.logger.Info("credential refresher stopped")
			return
		}
	}
}

func (rf *refresher) refreshAll() {
	ctx := context.Background()
	creds := rf.registry.GetAllCredentials()

	for _, cred := range creds {
		if cred.AuthType == "oauth2_client_credentials" && cred.NeedsRefresh() {
			rf.logger.Info("refreshing token",
				zap.Uint64("connectionID", cred.ConnectionID),
				zap.Time("expiresAt", cred.ExpiresAt),
			)

			if err := rf.registry.refreshToken(ctx, cred); err != nil {
				rf.logger.Error("failed to refresh token",
					zap.Uint64("connectionID", cred.ConnectionID),
					zap.Error(err),
				)
			} else {
				rf.logger.Info("token refreshed successfully",
					zap.Uint64("connectionID", cred.ConnectionID),
					zap.Time("newExpiresAt", cred.ExpiresAt),
				)
			}
		}
	}
}
