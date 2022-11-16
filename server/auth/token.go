package auth

import (
	"context"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	tokenService struct {
		store store.AuthOa2tokens
	}
)

func (svc tokenService) SearchByUserID(ctx context.Context, userID uint64) (types.AuthOa2tokenSet, error) {
	set, _, err := store.SearchAuthOa2tokens(ctx, svc.store, types.AuthOa2tokenFilter{UserID: userID})
	return set, err
}

func (svc tokenService) DeleteByID(ctx context.Context, ID uint64) error {
	return store.DeleteAuthOa2tokenByID(ctx, svc.store, ID)
}

func (svc tokenService) DeleteByUserID(ctx context.Context, userID uint64) error {
	return store.DeleteAuthOA2TokenByUserID(ctx, svc.store, userID)
}
