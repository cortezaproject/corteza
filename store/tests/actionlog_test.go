package tests

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
)

func testActionlog(t *testing.T, s store.Actionlogs) {
	var (
		ctx = context.Background()

		makeNew = func(id uint64, dd ...string) *actionlog.Action {
			// minimum data set for new user
			desc := strings.Join(dd, "")
			return &actionlog.Action{
				ID:          id,
				Timestamp:   *now(),
				Action:      "test-action",
				Resource:    "test-resource",
				Description: desc,
			}
		}

		stringifySetRange = func(set actionlog.ActionSet) string {
			if len(set) == 0 {
				return ""
			}

			var out = strconv.FormatUint(set[0].ID, 10)

			if len(set) > 1 {
				out += ".." + strconv.FormatUint(set[len(set)-1].ID, 10)
			}

			return out
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, actionlog.ActionSet) {
			req := require.New(t)
			req.NoError(s.TruncateActionlogs(ctx))

			set := make([]*actionlog.Action, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(uint64(i), string(rand.Bytes(10)))
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

	t.Run("with keyed paging", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateActionlogs(ctx))

		set := []*actionlog.Action{
			makeNew(1, "01"),
			makeNew(2, "02"),
			makeNew(3, "03"),
			makeNew(4, "04"),
			makeNew(5, "05"),
			makeNew(6, "06"),
			makeNew(7, "07"),
			makeNew(8, "08"),
			makeNew(9, "09"),
			makeNew(10, "10"),
		}

		req.NoError(s.CreateActionlog(ctx, set...))
		f := actionlog.Filter{}

		// Fetch first page
		f.Limit = 3
		set, f, err := store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 3)
		req.NotNil(f.NextPage)
		req.Nil(f.PrevPage)
		req.Equal("10..8", stringifySetRange(set))

		// 2nd page
		f.Limit = 6
		f.PageCursor = f.NextPage
		set, f, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 6)
		req.NotNil(f.NextPage)
		req.NotNil(f.PrevPage)
		req.Equal("7..2", stringifySetRange(set))

		// 3rd, last page (1 item left)
		f.Limit = 2
		f.PageCursor = f.NextPage
		set, f, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 1)
		req.NotNil(f.NextPage)
		req.NotNil(f.PrevPage)
		req.Equal("1", stringifySetRange(set))

		// try and go pass the last page
		f.PageCursor = f.NextPage
		set, _, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 0)

		// now, in reverse, last 3 items
		f.Limit = 3
		f.PageCursor = f.PrevPage
		set, f, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 3)
		req.NotNil(f.NextPage)
		req.NotNil(f.PrevPage)
		req.Equal("3..1", stringifySetRange(set))

		// still in reverse, next 6 items
		f.Limit = 5
		f.PageCursor = f.PrevPage
		set, f, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 5)
		req.NotNil(f.NextPage)
		req.NotNil(f.PrevPage)
		req.Equal("9..4", stringifySetRange(set))

		// still in reverse, last 5 items (actually, we'll only get 1)
		f.Limit = 5
		f.PageCursor = f.PrevPage
		set, f, err = store.SearchActionlogs(ctx, s, f)
		req.NoError(err)
		req.Len(set, 1)
		req.Nil(f.PrevPage)
		req.NotNil(f.NextPage)
		req.Equal("10", stringifySetRange(set))
	})
}
