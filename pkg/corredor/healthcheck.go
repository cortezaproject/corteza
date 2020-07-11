package corredor

import (
	"context"
	"fmt"
	"google.golang.org/grpc/connectivity"
)

// Healtcheck for global
func Healthcheck(ctx context.Context) error {
	if gCorredor == nil {
		return fmt.Errorf("uninitialized")
	}

	if !gCorredor.opt.Enabled {
		return nil
	}

	if state := gCorredor.conn.GetState(); state != connectivity.Ready {
		return fmt.Errorf("connection is %s", state)
	}

	return nil
}
