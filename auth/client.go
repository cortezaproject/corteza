package auth

import (
	"context"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	clientService struct {
		store interface {
			store.AuthClients
			store.AuthConfirmedClients
		}
	}
)

func (svc clientService) LookupByID(ctx context.Context, clientID uint64) (*types.AuthClient, error) {
	return store.LookupAuthClientByID(ctx, svc.store, clientID)
}

func (svc clientService) Confirmed(ctx context.Context, userID uint64) (types.AuthConfirmedClientSet, error) {
	set, _, err := store.SearchAuthConfirmedClients(ctx, svc.store, types.AuthConfirmedClientFilter{UserID: userID})
	return set, err
}

func (svc clientService) Revoke(ctx context.Context, userID, clientID uint64) error {
	return store.DeleteAuthConfirmedClientByUserIDClientID(ctx, svc.store, userID, clientID)
}
