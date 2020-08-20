package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testActionlog(t *testing.T, s store.Actionlogs) {
	var (
		ctx = context.Background()

		makeNew = func(dd ...string) *actionlog.Action {
			// minimum data set for new user
			desc := strings.Join(dd, "")
			return &actionlog.Action{
				ID:          id.Next(),
				Timestamp:   *now(),
				Action:      "test-action",
				Resource:    "test-resource",
				Description: desc,
			}
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, actionlog.ActionSet) {
			req := require.New(t)
			req.NoError(s.TruncateActionlogs(ctx))

			set := make([]*actionlog.Action, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateActionlog(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateActionlogs(ctx))
		action := &actionlog.Action{
			ID:        42,
			Timestamp: time.Now(),
			Resource:  "resource",
			Action:    "action",
		}
		req.NoError(s.CreateActionlog(ctx, action))
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by resource", func(t *testing.T) {
			req, set := truncAndFill(t, 5)

			set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Resource: "test-resource"})
			req.NoError(err)
			req.Len(set, 5)
		})

		t.Run("by action", func(t *testing.T) {
			req, set := truncAndFill(t, 5)

			set, _, err := s.SearchActionlogs(ctx, actionlog.Filter{Action: "test-action"})
			req.NoError(err)
			req.Len(set, 5)
		})
	})

}
