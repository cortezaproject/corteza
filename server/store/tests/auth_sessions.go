package tests

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func testAuthSessions(t *testing.T, s store.AuthSessions) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func(handle string) *types.AuthSession {
			// minimum data set for new authSession
			return &types.AuthSession{
				ID:        handle,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now(),
				Data:      []byte("..."),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.AuthSession) {
			req := require.New(t)
			req.NoError(s.TruncateAuthSessions(ctx))
			res := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateAuthSession(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		authSession := makeNew("AuthSessionCRUD")
		req.NoError(s.CreateAuthSession(ctx, authSession))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, authSession := truncAndCreate(t)
		fetched, err := s.LookupAuthSessionByID(ctx, authSession.ID)
		req.NoError(err)
		req.Equal(authSession.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
	})
	t.Run("delete", func(t *testing.T) {
		t.Run("by AuthSession", func(t *testing.T) {
			req, authSession := truncAndCreate(t)
			req.NoError(s.DeleteAuthSession(ctx, authSession))
			_, err := s.LookupAuthSessionByID(ctx, authSession.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})

		t.Run("by ID", func(t *testing.T) {
			req, authSession := truncAndCreate(t)
			req.NoError(s.DeleteAuthSessionByID(ctx, authSession.ID))
			_, err := s.LookupAuthSessionByID(ctx, authSession.ID)
			req.EqualError(err, store.ErrNotFound.Error())
		})
	})
}
