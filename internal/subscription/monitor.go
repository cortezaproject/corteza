package subscription

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/crusttech/crust/internal/logger"
)

// Starts subscription checker
func Monitor(ctx context.Context) context.Context {
	log := logger.Default()

	check := func(ctx context.Context) bool {
		log.Debug("validating subscription")
		if err := Check(ctx); err != nil {
			log.Error("subscription could not be validated", zap.Error(err))
			return false
		} else {
			log.Info("subscription validated")
		}

		return true
	}

	if !check(ctx) {
		// Initial subscription check failed,
		// Just exit.
		os.Exit(1)
	}

	// Initialize new context with cancellation we'll return this context and use it from this point on so that we make
	// a clean exist in case subscription checking fails
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		// Validate subscription key every 24h ours (from the last start of the service)
		t := time.NewTicker(time.Hour * 24)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				// Check the subscription again and call cancel on context
				if !check(ctx) {
					cancel()
					os.Exit(1)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return ctx
}
