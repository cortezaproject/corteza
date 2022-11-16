package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testDataPrivacyRequestComments(t *testing.T, s store.DataPrivacyRequestComments) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.DataPrivacyRequestComment {
			// minimum data set for new dataPrivacyRequestComment
			return &types.DataPrivacyRequestComment{
				ID:        id.Next(),
				CreatedAt: time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DataPrivacyRequestComment) {
			req := require.New(t)
			req.NoError(s.TruncateDataPrivacyRequestComments(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDataPrivacyRequestComment(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		dataPrivacyRequestComment := makeNew("DataPrivacyRequestCommentCRUD")
		req.NoError(s.CreateDataPrivacyRequestComment(ctx, dataPrivacyRequestComment))
	})

	t.Run("update", func(t *testing.T) {
		req, dataPrivacyRequestComment := truncAndCreate(t)
		dataPrivacyRequestComment.Comment = "DataPrivacyRequestCommentCRUD-2"

		req.NoError(s.UpdateDataPrivacyRequestComment(ctx, dataPrivacyRequestComment))
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, dataPrivacyRequestComment := truncAndCreate(t)
			dataPrivacyRequestComment.Comment = "DataPrivacyRequestCommentCRUD-2"

			req.NoError(s.UpsertDataPrivacyRequestComment(ctx, dataPrivacyRequestComment))
		})

		t.Run("new", func(t *testing.T) {
			dataPrivacyRequestComment := makeNew("upsert me")
			dataPrivacyRequestComment.Comment = "ComposeChartCRUD-2"

			req.NoError(s.UpsertDataPrivacyRequestComment(ctx, dataPrivacyRequestComment))
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by DataPrivacyRequestComment", func(t *testing.T) {
			req, dataPrivacyRequestComment := truncAndCreate(t)
			req.NoError(s.DeleteDataPrivacyRequestComment(ctx, dataPrivacyRequestComment))
		})

		t.Run("by ID", func(t *testing.T) {
			req, dataPrivacyRequestComment := truncAndCreate(t)
			req.NoError(s.DeleteDataPrivacyRequestCommentByID(ctx, dataPrivacyRequestComment.ID))
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.DataPrivacyRequestComment{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
			makeNew("/two-yes"),
		}

		count := len(prefill)
		valid := count

		req.NoError(s.TruncateDataPrivacyRequestComments(ctx))
		req.NoError(s.CreateDataPrivacyRequestComment(ctx, prefill...))

		// search for all valid
		set, _, err := s.SearchDataPrivacyRequestComments(ctx, types.DataPrivacyRequestCommentFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one
	})
}
