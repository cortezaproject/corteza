package app

import (
	"context"
	"fmt"
	authHandlers "github.com/cortezaproject/corteza-server/auth/handlers"
	"github.com/cortezaproject/corteza-server/auth/oauth2"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"io/ioutil"
	"strings"
)

func (app *CortezaApp) initAuth(ctx context.Context) (err error) {
	log := app.Log.Named("auth")

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

	// create token signature function from AUTH_ options
	// and pass it on to the token issuer
	alg, key, err := prepareSignatureFnParams(app.Opt.Auth)
	if err != nil {
		return fmt.Errorf("could not initialize token signer: %w", err)
	}

	log.Info(
		"initializing JWT and authentication procedures",
		zap.Stringer("algoritm", alg),
	)

	// construct token issuer with algorithm, secrets,
	auth.TokenIssuer, err = auth.NewTokenIssuer(
		auth.WithSigner(func(t jwt.Token) ([]byte, error) {
			return jwt.Sign(t, alg, key)
		}),

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

	return
}

// helper function that loads and/or parses private keys and initializes JWT signer function
// from the given arguments
func prepareSignatureFnParams(opt options.AuthOpt) (alg jwa.SignatureAlgorithm, key interface{}, err error) {
	alg = jwa.SignatureAlgorithm(opt.JwtAlgorithm)

	switch alg {
	case jwa.HS256, jwa.HS384, jwa.HS512:
		// expecting secret to be set
		if len(opt.Secret) == 0 {
			return alg, nil, fmt.Errorf("token secret missing")
		}

		key = []byte(opt.Secret)
	case
		jwa.PS256, jwa.PS384, jwa.PS512,
		jwa.RS256, jwa.RS384, jwa.RS512:
		if len(opt.JwtKey) == 0 {
			return alg, nil, fmt.Errorf("token key missing")
		}

		// if given key dos not begins with "-----BEGIN"
		// assume it's path to a file and load contents of that file
		if !strings.HasPrefix(opt.JwtKey, "-----BEGIN") {
			var (
				keyFile = opt.JwtKey
				b       []byte
			)
			if b, err = ioutil.ReadFile(keyFile); err != nil {
				return alg, nil, fmt.Errorf("could not read key file: %w", err)
			}

			// overwrite th input and load
			opt.JwtKey = string(b)

			// recheck contents of the key
			if !strings.HasPrefix(opt.JwtKey, "-----BEGIN") {
				return alg, nil, fmt.Errorf("file %q does not contain a valid private key", keyFile)
			}
		}

		// generates pem.Private from the kInput string
		key, err = jwk.ParseKey([]byte(opt.JwtKey), jwk.WithPEM(true))
	case "":
		// should be caught by options init procedure and set to default,
		// but you never know...
		err = fmt.Errorf("token signature algorithm empty or missing")

	default:
		err = fmt.Errorf("token signature algorithm %q not supported", alg)
	}

	return
}
