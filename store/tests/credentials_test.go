package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testCredentials(t *testing.T, s store.Credentials) {
	var (
		ctx     = context.Background()
		makeNew = func(nn ...string) *types.Credentials {
			// minimum data set for new user
			name := strings.Join(nn, "")
			return &types.Credentials{
				ID:          id.Next(),
				CreatedAt:   time.Now(),
				Credentials: name,
				Kind:        "test-kind",
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Credentials) {
			req := require.New(t)
			req.NoError(s.TruncateCredentials(ctx))
			user := makeNew()
			req.NoError(s.CreateCredentials(ctx, user))
			return req, user
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateCredentials(ctx))

		credentials := &types.Credentials{
			ID:        42,
			CreatedAt: time.Now(),
			Label:     "CredentialsCRUD",
		}
		req.NoError(s.CreateCredentials(ctx, credentials))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		req, crd := truncAndCreate(t)
		fetched, err := s.LookupCredentialsByID(ctx, crd.ID)
		req.NoError(err)
		req.Equal(crd.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {
		req, crd := truncAndCreate(t)
		crd.Credentials = "new-credentials"
		req.NoError(s.UpdateCredentials(ctx, crd))
		fetched, err := s.LookupCredentialsByID(ctx, crd.ID)
		req.NoError(err)
		req.Equal("new-credentials", fetched.Credentials)

	})
}
