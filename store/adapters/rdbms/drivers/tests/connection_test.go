package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestConnectionPing(t *testing.T) {
	var (
		ctx = context.Background()
	)

	eachDB(t, func(t *testing.T, c *conn) error {
		toCtx, cfn := context.WithTimeout(ctx, time.Second)
		defer cfn()

		if err := c.db.PingContext(toCtx); err != nil {
			return fmt.Errorf("can not ping database: %w", err)
		}

		return nil
	})
}
