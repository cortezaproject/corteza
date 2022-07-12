package app

import (
	"context"
	"fmt"
	authHandlers "github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/oauth2"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

func (app *CortezaApp) initAuth(ctx context.Context) (err error) {
	if app.Opt.Auth.DefaultClient != "" {
		// default client will help streamline authorization with default clients
		app.DefaultAuthClient, err = store.LookupAuthClientByHandle(ctx, app.Store, app.Opt.Auth.DefaultClient)
		if err != nil {
			return fmt.Errorf("cannot load default client: %w", err)
		}
	}

	app.oa2m = oauth2.NewManager(
		app.Opt.Auth,
		app.Log,
		oauth2.NewClientStore(app.Store, app.DefaultAuthClient),
		oauth2.NewTokenStore(app.Store),
	)

	// set base path for links&routes in auth server
	authHandlers.BasePath = app.Opt.HTTPServer.BaseUrl

	auth.DefaultSigner = auth.HmacSigner(app.Opt.Auth.Secret)

	if auth.HttpTokenVerifier, err = auth.TokenVerifierMiddlewareWithSecretSigner(app.Opt.Auth.Secret); err != nil {
		return fmt.Errorf("could not set token verifier")
	}

	auth.TokenIssuer, err = auth.NewTokenIssuer(
		auth.WithSecretSigner(app.Opt.Auth.Secret),
		// @todo implement configurable issuer claim
		//auth.WithDefaultIssuer(app.Opt.Auth.TokenClaimIssuer),
		auth.WithDefaultExpiration(app.Opt.Auth.AccessTokenLifetime),
		auth.WithDefaultClientID(app.DefaultAuthClient.ID),
		auth.WithLookup(func(ctx context.Context, accessToken string) (err error) {
			_, err = store.LookupAuthOa2tokenByAccess(ctx, app.Store, accessToken)
			return err
		}),
		auth.WithStore(func(ctx context.Context, req auth.TokenRequest) error {
			var (
				eti       = auth.GetExtraReqInfoFromContext(ctx)
				createdAt = req.IssuedAt

				oa2t = &types.AuthOa2token{
					ID:         id.Next(),
					Access:     req.AccessToken,
					Refresh:    req.RefreshToken,
					CreatedAt:  createdAt,
					RemoteAddr: eti.RemoteAddr,
					UserAgent:  eti.UserAgent,
					ClientID:   req.ClientID,
					UserID:     req.UserID,
					ExpiresAt:  createdAt.Add(req.Expiration),
				}
			)

			return store.CreateAuthOa2token(ctx, app.Store, oa2t)
		}),
	)

	if err != nil {
		return fmt.Errorf("could not initialize token issuer: %w", err)
	}
}
