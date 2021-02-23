package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-oauth2/oauth2/v4"
	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	oauth2models "github.com/go-oauth2/oauth2/v4/models"
	"strconv"
	"time"
)

type (
	CortezaTokenStore struct {
		Store interface {
			store.AuthOa2tokens
			store.AuthConfirmedClients
		}
	}
)

var (
	// wrapper around time.Now() that will aid service testing
	now = func() *time.Time {
		c := time.Now().Round(time.Second)
		return &c
	}

	// wrapper around id.Next() that will aid service testing
	nextID = func() uint64 {
		return id.Next()
	}

	_ oauth2.TokenStore = &CortezaTokenStore{}
)

func (c CortezaTokenStore) Create(ctx context.Context, info oauth2.TokenInfo) (err error) {
	var (
		eti  = request.GetExtraReqInfoFromContext(ctx)
		oa2t = &types.AuthOa2token{
			ID:         nextID(),
			CreatedAt:  *now(),
			RemoteAddr: eti.RemoteAddr,
			UserAgent:  eti.UserAgent,
		}

		acc = &types.AuthConfirmedClient{
			ConfirmedAt: oa2t.CreatedAt,
		}
	)

	if code := info.GetCode(); code != "" {
		oa2t.Code = code
		oa2t.ExpiresAt = info.GetCodeCreateAt().Add(info.GetCodeExpiresIn())
	} else {
		oa2t.Access = info.GetAccess()
		oa2t.ExpiresAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())

		if refresh := info.GetRefresh(); refresh != "" {
			oa2t.Refresh = info.GetRefresh()
			oa2t.ExpiresAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}

	if oa2t.Data, err = json.Marshal(info); err != nil {
		return
	}

	if oa2t.ClientID, err = strconv.ParseUint(info.GetClientID(), 10, 64); err != nil {
		return fmt.Errorf("could not parse client ID from token info: %w", err)
	}

	// copy client id to auth client confirmation
	acc.ClientID = oa2t.ClientID

	if info.GetUserID() != "" {
		if oa2t.UserID = auth.ExtractUserIDFromSubClaim(info.GetUserID()); oa2t.UserID == 0 {
			// UserID stores collection of IDs: user's ID and set of all roles user is member of
			return fmt.Errorf("could not parse user ID from token info: %w", err)
		}
	}

	// copy user id to auth client confirmation
	acc.UserID = oa2t.UserID

	if err = store.UpsertAuthConfirmedClient(ctx, c.Store, acc); err != nil {
		return
	}

	return store.CreateAuthOa2token(ctx, c.Store, oa2t)
}

func (c CortezaTokenStore) RemoveByCode(ctx context.Context, code string) error {
	return store.DeleteAuthOA2TokenByCode(ctx, c.Store, code)
}

func (c CortezaTokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return store.DeleteAuthOA2TokenByAccess(ctx, c.Store, access)
}

func (c CortezaTokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return store.DeleteAuthOA2TokenByRefresh(ctx, c.Store, refresh)
}

func (c CortezaTokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	var (
		internal = &oauth2models.Token{}
		t, err   = store.LookupAuthOa2tokenByCode(ctx, c.Store, code)
	)
	if err != nil {

		if errors.IsNotFound(err) {
			return nil, oauth2errors.ErrInvalidAuthorizeCode
		}

		return nil, fmt.Errorf("failed to get code: %w", err)
	}

	return internal, t.Data.Unmarshal(internal)
}

func (c CortezaTokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	var (
		internal = &oauth2models.Token{}
		t, err   = store.LookupAuthOa2tokenByAccess(ctx, c.Store, access)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return internal, t.Data.Unmarshal(internal)
}

func (c CortezaTokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	var (
		internal = &oauth2models.Token{}
		t, err   = store.LookupAuthOa2tokenByRefresh(ctx, c.Store, refresh)
	)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, oauth2errors.ErrInvalidRefreshToken
		}

		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if t.ExpiresAt.Before(*now()) {
		//return nil, oauth2errors.ErrExpiredRefreshToken
	}

	return internal, t.Data.Unmarshal(internal)
}
