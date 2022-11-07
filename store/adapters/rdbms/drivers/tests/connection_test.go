package tests

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConnectionPing(t *testing.T) {
	t.Logf("testing %s", conn.dsn)

	toCtx, cfn := context.WithTimeout(ctx, time.Second)
	defer cfn()

	require.NoError(t, conn.db.PingContext(toCtx))
}
