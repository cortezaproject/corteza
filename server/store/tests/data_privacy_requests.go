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

func testDataPrivacyRequests(t *testing.T, s store.DataPrivacyRequests) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.DataPrivacyRequest {
			// minimum data set for new dataPrivacyRequest
			return &types.DataPrivacyRequest{
				ID:          id.Next(),
				CreatedAt:   time.Now(),
				RequestedAt: time.Now(),
				Kind:        types.RequestKindCorrect,
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.DataPrivacyRequest) {
			req := require.New(t)
			req.NoError(s.TruncateDataPrivacyRequests(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateDataPrivacyRequest(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		dataPrivacyRequest := makeNew("DataPrivacyRequestCRUD")
		req.NoError(s.CreateDataPrivacyRequest(ctx, dataPrivacyRequest))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, dataPrivacyRequest := truncAndCreate(t)
		fetched, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
		req.NoError(err)
		req.Equal(dataPrivacyRequest.Kind, fetched.Kind)
		req.Equal(dataPrivacyRequest.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, dataPrivacyRequest := truncAndCreate(t)
		dataPrivacyRequest.Kind = types.RequestKindDelete

		req.NoError(s.UpdateDataPrivacyRequest(ctx, dataPrivacyRequest))

		updated, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
		req.NoError(err)
		req.Equal(dataPrivacyRequest.Kind, updated.Kind)
	})

	t.Run("upsert", func(t *testing.T) {
		t.Run("existing", func(t *testing.T) {
			req, dataPrivacyRequest := truncAndCreate(t)
			dataPrivacyRequest.Kind = types.RequestKindDelete

			req.NoError(s.UpsertDataPrivacyRequest(ctx, dataPrivacyRequest))

			updated, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
			req.NoError(err)
			req.Equal(dataPrivacyRequest.Kind, updated.Kind)
		})

		t.Run("new", func(t *testing.T) {
			dataPrivacyRequest := makeNew("upsert me")
			dataPrivacyRequest.Kind = "ComposeChartCRUD-2"

			req.NoError(s.UpsertDataPrivacyRequest(ctx, dataPrivacyRequest))

			upserted, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
			req.NoError(err)
			req.Equal(dataPrivacyRequest.Kind, upserted.Kind)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by DataPrivacyRequest", func(t *testing.T) {
			req, dataPrivacyRequest := truncAndCreate(t)
			req.NoError(s.DeleteDataPrivacyRequest(ctx, dataPrivacyRequest))
			_, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, dataPrivacyRequest := truncAndCreate(t)
			req.NoError(s.DeleteDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID))
			_, err := s.LookupDataPrivacyRequestByID(ctx, dataPrivacyRequest.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})

	t.Run("search", func(t *testing.T) {
		prefill := []*types.DataPrivacyRequest{
			makeNew("/one-one"),
			makeNew("/one-two"),
			makeNew("/two-one"),
			makeNew("/two-two"),
		}

		prefill[3].Kind = types.RequestKindDelete

		count := len(prefill)
		valid := count

		req.NoError(s.TruncateDataPrivacyRequests(ctx))
		req.NoError(s.CreateDataPrivacyRequest(ctx, prefill...))

		// search for all valid
		set, _, err := s.SearchDataPrivacyRequests(ctx, types.DataPrivacyRequestFilter{})
		req.NoError(err)
		req.Len(set, valid) // we've deleted one

		set, _, err = s.SearchDataPrivacyRequests(ctx, types.DataPrivacyRequestFilter{Kind: []string{"delete"}})
		req.NoError(err)
		req.Len(set, 1)
	})
}
