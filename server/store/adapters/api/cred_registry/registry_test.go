package cred_registry

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockStore struct{}

func (m *mockStore) LookupDalConnectionByID(ctx context.Context, id uint64) (*types.DalConnection, error) {
	return &types.DalConnection{
		ID: id,
		Config: types.ConnectionConfig{
			DAL: &types.ConnectionConfigDAL{
				Params: map[string]any{
					"auth": map[string]any{
						"method": "oauth2_client_credentials",
						"params": map[string]any{},
					},
				},
			},
		},
	}, nil
}

func (m *mockStore) UpdateDalConnection(ctx context.Context, conn ...*types.DalConnection) error {
	return nil
}

func TestRegistry_StoreAndGet(t *testing.T) {
	reg, err := New(&mockStore{}, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, reg)

	cred := &Credential{
		ConnectionID: 1,
		AuthType:     "bearer",
		Token:        "test-token",
	}

	err = reg.Store(cred)
	require.NoError(t, err)

	retrieved, err := reg.Get(1)
	require.NoError(t, err)
	require.Equal(t, "test-token", retrieved.Token)
}

func TestCredential_NeedsRefresh(t *testing.T) {
	cred := &Credential{
		AuthType:  "oauth2_client_credentials",
		ExpiresAt: time.Now().Add(3 * time.Minute), // Expires in 3 minutes
	}

	// Should need refresh (< 5 min buffer)
	require.True(t, cred.NeedsRefresh())

	cred.ExpiresAt = time.Now().Add(10 * time.Minute)
	// Should not need refresh (> 5 min buffer)
	require.False(t, cred.NeedsRefresh())
}
