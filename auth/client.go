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

func (svc clientService) Lookup(ctx context.Context, identifier interface{}) (*types.AuthClient, error) {
	return clientLookup(ctx, svc.store, identifier)
}

func (svc clientService) LookupByHandle(ctx context.Context, handle string) (*types.AuthClient, error) {
	return store.LookupAuthClientByHandle(ctx, svc.store, handle)
}

func (svc clientService) Confirmed(ctx context.Context, userID uint64) (types.AuthConfirmedClientSet, error) {
	set, _, err := store.SearchAuthConfirmedClients(ctx, svc.store, types.AuthConfirmedClientFilter{UserID: userID})
	return set, err
}

func (svc clientService) Revoke(ctx context.Context, userID, clientID uint64) error {
	return store.DeleteAuthConfirmedClientByUserIDClientID(ctx, svc.store, userID, clientID)
}
