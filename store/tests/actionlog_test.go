package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testActionlog(t *testing.T, s actionlogsStore) {
	var (
		ctx = context.Background()
		req = require.New(t)

		//err  error
		action *actionlog.Action
	)

	t.Run("create", func(t *testing.T) {
		action = &actionlog.Action{
			ID:        42,
			Timestamp: time.Now(),
			Resource:  "resource",
			Action:    "action",
		}
		req.NoError(s.CreateActionlog(ctx, action))
	})

	t.Run("search by resource", func(t *testing.T) {
		set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Resource: action.Resource})
		req.NoError(err)
		req.Len(set, 1)
	})

	t.Run("search by action", func(t *testing.T) {
		set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Action: action.Action})
		req.NoError(err)
		req.Len(set, 1)
	})

}
