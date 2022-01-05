package oauth2

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"go.uber.org/zap"
)

const (
	RedirectUriSeparator = " "
)

func NewManager(opt options.AuthOpt, log *zap.Logger, cs oauth2.ClientStore, ts oauth2.TokenStore) *manage.Manager {
	manager := manage.NewDefaultManager()

	// Here we are cloning the internal package variable as I do not think
	// it is sane to overwrite it directly.
	cfg := *manage.DefaultAuthorizeCodeTokenCfg
	cfg.AccessTokenExp = opt.AccessTokenLifetime
	cfg.RefreshTokenExp = opt.RefreshTokenLifetime

	manager.SetAuthorizeCodeTokenCfg(&cfg)

	// token store
	manager.MapTokenStorage(ts)

	// generate jwt access token
	manager.MapAccessGenerate(NewJWTAccessGenerate(auth.DefaultJwtHandler))
	manager.MapClientStorage(cs)

	manager.SetValidateURIHandler(func(baseURI, redirectURI string) (err error) {
		if baseURI == "" {
			log.Debug(
				"redirect URI check for client is disabled (empty validation list)",
				zap.String("sent", redirectURI),
			)

			return nil
		}

		var (
			valid = strings.Split(baseURI, RedirectUriSeparator)
		)

		log.Debug(
			"matching redirectURI",
			zap.String("sent", redirectURI),
			zap.Strings("valid", valid),
		)

		for _, baseURI = range valid {
			if strings.HasPrefix(redirectURI, baseURI) {
				return nil
			}
		}

		return errors.ErrInvalidRedirectURI
	})

	return manager
}

func NewServer(manager *manage.Manager) *server.Server {
	srv := server.NewServer(&server.Config{
		TokenType:             "Bearer",
		AllowGetAccessRequest: false,
		AllowedResponseTypes: []oauth2.ResponseType{
			oauth2.Code,
		},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.Refreshing,
			// before enabling ClientCredentials grant type, we need to know how to modify released token
			// using client's security info; how to enforce impersonated user and his roles.
			oauth2.ClientCredentials,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengePlain,
			oauth2.CodeChallengeS256,
		},
	}, manager)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		return errors.NewResponse(err, 500)
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		msg := re.Description
		if msg == "" {
			msg = re.Error.Error()
		}

		logger.Default().
			Warn(msg)
	})

	return srv
}
