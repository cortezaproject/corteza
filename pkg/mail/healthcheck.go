package mail

import (
	"context"
	"fmt"
)

// Healtcheck for (global) scheduler
func Healthcheck(ctx context.Context) error {
	if defaultDialer == nil {
		return fmt.Errorf("uninitialized")
	}

	if defaultDialerError != nil {
		return defaultDialerError
	}

	return nil
}
