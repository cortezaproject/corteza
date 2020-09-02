package mail

import (
	"context"
	"fmt"
)

// Healtcheck for (global) scheduler
func Healthcheck(ctx context.Context) error {
	if defaultDialerError == nil {
		return defaultDialerError
	}

	if defaultDialer == nil {
		return fmt.Errorf("uninitialized")
	}

	return nil
}
