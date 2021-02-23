package oauth2

import (
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

const (
	RedirectUriSeparator = " "
)

func NewManager(opt options.AuthOpt, cs oauth2.ClientStore, ts oauth2.TokenStore) *manage.Manager {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	manager.MapTokenStorage(ts)

	// generate jwt access token
	manager.MapAccessGenerate(NewJWTAccessGenerate("", []byte(opt.Secret), jwt.SigningMethodHS512))
	manager.MapClientStorage(cs)

	manager.SetValidateURIHandler(func(baseURI, redirectURI string) (err error) {
		var (
			base, redirect *url.URL
		)

		redirect, err = url.Parse(redirectURI)
		if err != nil {
			return err
		}

		// allow port only when using localhost as redirect
		if redirect.Port() != "" && redirect.Hostname() != "localhost" {
			return errors.ErrInvalidRedirectURI
		}

		for _, baseURI = range strings.Split(baseURI, RedirectUriSeparator) {
			base, err = url.Parse(baseURI)
			if err != nil {
				return err
			}

			if strings.HasPrefix(redirect.String(), base.String()) {
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
			//
			// oauth2.ClientCredentials,
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
			WithOptions(zap.AddStacktrace(zap.PanicLevel)).
			Error(msg)
	})

	return srv
}
