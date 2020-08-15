package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/system/types"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func testCredentials(t *testing.T, s credentialsStore) {
	var (
		ctx         = context.Background()
		req         = require.New(t)
		credentials *types.Credentials
		//err         error
	)

	t.Run("create", func(t *testing.T) {
		credentials = &types.Credentials{
			ID:        42,
			CreatedAt: time.Now(),
			Label:     "CredentialsCRUD",
		}
		req.NoError(s.CreateCredentials(ctx, credentials))
	})

	t.Run("lookup by ID", func(t *testing.T) {
		fetched, err := s.LookupCredentialsByID(ctx, credentials.ID)
		req.NoError(err)
		req.Equal(credentials.ID, fetched.ID)
		req.NotNil(fetched.CreatedAt)
		req.Nil(fetched.UpdatedAt)
		req.Nil(fetched.DeletedAt)
	})

	t.Run("update", func(t *testing.T) {

		credentials = &types.Credentials{
			ID:        42,
			CreatedAt: time.Now(),
			Label:     "CredentialsCRUD+2",
		}
		req.NoError(s.UpdateCredentials(ctx, credentials))
	})

	//t.Run("delete/undelete", func(t *testing.T) {
	//	ID := credentials.ID
	//	credentials, err = s.LookupCredentialsByID(ctx, ID)
	//	req.NoError(err)
	//
	//	req.NoError(s.DeleteCredentialsByID(ctx, ID))
	//	credentials, err = s.LookupCredentialsByID(ctx, ID)
	//	req.NoError(err)
	//	req.NotNil(credentials.DeletedAt)
	//
	//	req.NoError(s.UndeleteCredentialsByID(ctx, ID))
	//	credentials, err = s.LookupCredentialsByID(ctx, ID)
	//	req.NoError(err)
	//	req.Nil(credentials.DeletedAt)
	//})

	t.Run("search by *", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
