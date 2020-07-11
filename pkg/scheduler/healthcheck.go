package scheduler

import (
	"context"
	"fmt"
)

// Healtcheck for (global) scheduler
func Healthcheck(ctx context.Context) error {
	if gScheduler == nil {
		return fmt.Errorf("uninitialized")
	}

	if gScheduler.ticker == nil {
		return fmt.Errorf("stopped")
	}

	return nil
}
