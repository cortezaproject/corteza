package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-oauth2/oauth2/v4"
	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	oauth2models "github.com/go-oauth2/oauth2/v4/models"
)

type (
	tokenStorer interface {
		store.AuthOa2tokens
		store.AuthConfirmedClients
	}

	tokenStore struct {
		Store tokenStorer
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

	_ oauth2.TokenStore = &tokenStore{}
)

func NewTokenStore(s tokenStorer) *tokenStore {
	return &tokenStore{Store: s}
}

func (c tokenStore) Create(ctx context.Context, info oauth2.TokenInfo) (err error) {
	var (
		oa2t *types.AuthOa2token
		acc  *types.AuthConfirmedClient

		userID   uint64
		clientID uint64
	)

	if clientID, err = strconv.ParseUint(info.GetClientID(), 10, 64); err != nil {
		return fmt.Errorf("could not parse client ID from token info: %w", err)
	}

	if userID, _ = auth.ExtractFromSubClaim(info.GetUserID()); userID == 0 {
		return fmt.Errorf("could not parse user ID from token info")
	}

	// Make oauth2 token and auth confirmation structs from user and client IDs
	if oa2t, acc, err = makeAuthStructs(ctx, userID, clientID, info, info.GetCodeExpiresIn()); err != nil {
		return
	}

	if err = store.UpsertAuthConfirmedClient(ctx, c.Store, acc); err != nil {
		return
	}

	if err = store.CreateAuthOa2token(ctx, c.Store, oa2t); err != nil {
		return
	}

	return nil
}

func (c tokenStore) RemoveByCode(ctx context.Context, code string) error {
	return store.DeleteAuthOA2TokenByCode(ctx, c.Store, code)
}

func (c tokenStore) RemoveByAccess(ctx context.Context, access string) error {
	return store.DeleteAuthOA2TokenByAccess(ctx, c.Store, access)
}

func (c tokenStore) RemoveByRefresh(ctx context.Context, refresh string) error {
	return store.DeleteAuthOA2TokenByRefresh(ctx, c.Store, refresh)
}

func (c tokenStore) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
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

func (c tokenStore) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	var (
		internal = &oauth2models.Token{}
		t, err   = store.LookupAuthOa2tokenByAccess(ctx, c.Store, access)
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return internal, t.Data.Unmarshal(internal)
}

func (c tokenStore) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
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

	return internal, t.Data.Unmarshal(internal)
}

func makeAuthStructs(ctx context.Context, userID, clientID uint64, info oauth2.TokenInfo, expiresAt time.Duration) (oa2t *types.AuthOa2token, acc *types.AuthConfirmedClient, err error) {
	var (
		eti       = auth.GetExtraReqInfoFromContext(ctx)
		createdAt = time.Now().Round(time.Second)
	)

	oa2t = &types.AuthOa2token{
		ID:         id.Next(),
		CreatedAt:  createdAt,
		RemoteAddr: eti.RemoteAddr,
		UserAgent:  eti.UserAgent,
		ClientID:   clientID,
		UserID:     userID,
		ExpiresAt:  createdAt.Add(expiresAt),
	}

	acc = &types.AuthConfirmedClient{
		ClientID:    clientID,
		UserID:      userID,
		ConfirmedAt: createdAt,
	}

	if oa2t.Data, err = json.Marshal(info); err != nil {
		return nil, nil, err
	}

	if code := info.GetCode(); code != "" {
		oa2t.Code = code
	} else {
		// When creating non-access-code tokens,
		// we need to overwrite expiration time
		// with custom values for access or refresh token
		oa2t.Access = info.GetAccess()
		oa2t.ExpiresAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())

		if refresh := info.GetRefresh(); refresh != "" {
			oa2t.Refresh = info.GetRefresh()
			oa2t.ExpiresAt = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		}
	}

	return
}
