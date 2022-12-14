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

func testAuthOa2tokens(t *testing.T, s store.AuthOa2tokens) {
	var (
		ctx = context.Background()
		req = require.New(t)

		makeNew = func() *types.AuthOa2token {
			// minimum data set for new authOa2token
			return &types.AuthOa2token{
				ID:        id.Next(),
				CreatedAt: time.Now(),
				ExpiresAt: time.Now(),
				Access:    string(rand.Bytes(5)),
				Code:      string(rand.Bytes(5)),
				Refresh:   string(rand.Bytes(5)),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.AuthOa2token) {
			req := require.New(t)
			req.NoError(s.TruncateAuthOa2tokens(ctx))
			res := makeNew()
			req.NoError(s.CreateAuthOa2token(ctx, res))
			return req, res
		}
	)

	t.Run("create", func(t *testing.T) {
		authOa2token := makeNew()
		req.NoError(s.CreateAuthOa2token(ctx, authOa2token))
	})

	t.Run("lookup by", func(t *testing.T) {
		t.Run("access", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			fetched, err := s.LookupAuthOa2tokenByAccess(ctx, authOa2token.Access)
			req.NoError(err)
			req.NotNil(fetched)
		})

		t.Run("code", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			fetched, err := s.LookupAuthOa2tokenByCode(ctx, authOa2token.Code)
			req.NoError(err)
			req.NotNil(fetched)
		})

		t.Run("refresh", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			fetched, err := s.LookupAuthOa2tokenByRefresh(ctx, authOa2token.Refresh)
			req.NoError(err)
			req.NotNil(fetched)
		})
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("access", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			err := s.DeleteAuthOA2TokenByRefresh(ctx, authOa2token.Access)
			req.NoError(err)
		})

		t.Run("code", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			err := s.DeleteAuthOA2TokenByRefresh(ctx, authOa2token.Code)
			req.NoError(err)
		})

		t.Run("refresh", func(t *testing.T) {
			req, authOa2token := truncAndCreate(t)
			err := s.DeleteAuthOA2TokenByRefresh(ctx, authOa2token.Refresh)
			req.NoError(err)
		})
	})
}
