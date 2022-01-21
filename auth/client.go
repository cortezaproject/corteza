package auth

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/store"
	systemService "github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
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

func clientLookup(ctx context.Context, s store.AuthClients, identifier interface{}) (*types.AuthClient, error) {
	if id := cast.ToUint64(identifier); id > 0 {
		return store.LookupAuthClientByID(ctx, s, id)
	} else if h := cast.ToString(identifier); handle.IsValid(h) {
		return store.LookupAuthClientByHandle(ctx, s, h)
	} else {
		return nil, systemService.AuthClientErrInvalidID()
	}
}
