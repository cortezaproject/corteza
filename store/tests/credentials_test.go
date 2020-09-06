package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	// "github.com/cortezaproject/corteza-server/pkg/rh"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func testCredentials(t *testing.T, s store.Credentials) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.Credentials {
			name := strings.Join(nn, "")
			return &types.Credentials{
				ID:          id.Next(),
				OwnerID:     id.Next(),
				Kind:        "test-kind" + name,
				Credentials: name,
				Label:       "CredentialsCRUD" + name,
				CreatedAt:   time.Now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.Credentials) {
			req := require.New(t)
			req.NoError(s.TruncateCredentials(ctx))
			credentials := makeNew()
			req.NoError(s.CreateCredentials(ctx, credentials))
			return req, credentials
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.CredentialsSet) {
			req := require.New(t)
			req.NoError(s.TruncateCredentials(ctx))

			set := make([]*types.Credentials, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateCredentials(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateCredentials(ctx, makeNew()))
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

	t.Run("update", func(t *testing.T) {
		req, crd := truncAndCreate(t)
		crd.Credentials = "new-credentials"
		req.NoError(s.UpdateCredentials(ctx, crd))
		fetched, err := s.LookupCredentialsByID(ctx, crd.ID)
		req.NoError(err)
		req.Equal("new-credentials", fetched.Credentials)
	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by credentials", func(t *testing.T) {
			req, crd := truncAndCreate(t)
			req.NoError(s.DeleteCredentials(ctx, crd))
			set, _, err := s.SearchCredentials(ctx, types.CredentialsFilter{OwnerID: crd.OwnerID})
			req.NoError(err)
			req.Len(set, 0)
		})

		t.Run("by ID", func(t *testing.T) {
			req, crd := truncAndCreate(t)
			req.NoError(s.DeleteCredentialsByID(ctx, crd.ID))
			set, _, err := s.SearchCredentials(ctx, types.CredentialsFilter{OwnerID: crd.OwnerID})
			req.NoError(err)
			req.Len(set, 0)
		})
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by owner", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchCredentials(ctx, types.CredentialsFilter{OwnerID: prefill[0].OwnerID})
			req.NoError(err)
			req.Len(set, 1)
		})

		t.Run("by kind", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchCredentials(ctx, types.CredentialsFilter{Kind: prefill[0].Kind})
			req.NoError(err)
			req.Len(set, 1)
		})

		// Not sure why this doesnt work
		// t.Run("by state", func(t *testing.T) {
		// 	t.Run("deleted", func(t *testing.T) {
		// 		req, prefill := truncAndFill(t, 5)
		// 		time := time.Now()
		// 		prefill[0].DeletedAt = &time
		// 		req.NoError(s.DeleteCredentialsByID(ctx, prefill[0].ID))

		// 		set, _, err := s.SearchCredentials(ctx, types.CredentialsFilter{Deleted: filter.StateExcluded})
		// 		req.NoError(err)
		// 		req.Len(set, 4)

		// 		set, _, err = s.SearchCredentials(ctx, types.CredentialsFilter{Deleted: filter.StateInclusive})
		// 		req.NoError(err)
		// 		req.Len(set, 5)

		// 		set, _, err = s.SearchCredentials(ctx, types.CredentialsFilter{Deleted: filter.StateExclusive})
		// 		req.NoError(err)
		// 		req.Len(set, 1)
		// 	})
		// })
	})
}
