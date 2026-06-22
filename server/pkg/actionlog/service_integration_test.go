package actionlog_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/actionlog"
	"github.com/crusttech/human/server/store"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/mysql"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/postgres"
	_ "github.com/crusttech/human/server/store/adapters/rdbms/drivers/sqlite"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupStore(t *testing.T) store.Storer {
	t.Helper()

	dsn, ok := os.LookupEnv("ACTIONLOG_DB_DSN")
	if !ok || dsn == "" {
		t.Skip("ACTIONLOG_DB_DSN not set")
	}

	ctx := context.Background()
	s, err := store.Connect(ctx, zap.NewNop(), dsn, false)
	require.NoError(t, err)
	require.NoError(t, store.Upgrade(ctx, zap.NewNop(), s))

	return s
}

func TestActionlog(t *testing.T) {
	s := setupStore(t)
	ctx := context.Background()
	req := require.New(t)

	svc := actionlog.NewService(s, zap.NewNop(), zap.NewNop(), actionlog.MakeDebugPolicy())

	req.NoError(store.TruncateActionlogs(ctx, s))

	svc.Record(ctx, &actionlog.Action{
		Timestamp: time.Now(),
		Resource:  "corteza::compose:module/1/42",
		Action:    "lookup",
		Severity:  actionlog.Notice,
	})

	set, _, err := svc.Find(ctx, actionlog.Filter{
		Resource: "corteza::compose:module/1/42",
	})
	req.NoError(err)
	req.Len(set, 1)
	req.Equal("corteza::compose:module/1/42", set[0].Resource)
}
