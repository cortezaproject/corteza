package oauth2

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/handle"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/pkg/payload"
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
	manager.MapClientStorage(cs)
	// Change the default config for it to update refresh token timestamps
	// else the refresh token timestamp remains the same
	//
	// @note do this so we don't change the default `manage` package var
	rcfg := *manage.DefaultRefreshTokenCfg
	rcfg.IsResetRefreshTime = true
	manager.SetRefreshTokenCfg(&rcfg)

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

func NewServer(manager oauth2.Manager) *server.Server {
	srv := server.NewServer(&server.Config{
		TokenType:             "Bearer",
		AllowGetAccessRequest: false,
		AllowedResponseTypes: []oauth2.ResponseType{
			oauth2.Code,
		},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.Refreshing,
			oauth2.ClientCredentials,
		},
		AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
			oauth2.CodeChallengePlain,
			oauth2.CodeChallengeS256,
		},
	}, manager)

	srv.ClientInfoHandler = func(r *http.Request) (clientID, clientSecret string, err error) {
		// check in basic handler first
		clientID, clientSecret, err = server.ClientBasicHandler(r)

		if clientID == "" && clientSecret == "" {
			//error or no error, if ID & secret are empty,
			// check the form handler
			clientID, clientSecret, err = server.ClientFormHandler(r)
		}

		// just in case, when client's handle is used instead of the ID
		// preload it here
		if id := payload.ParseUint64(clientID); id == 0 && handle.IsValid(clientID) {
			var client oauth2.ClientInfo
			client, err = manager.GetClient(r.Context(), clientID)
			if err != nil {
				err = fmt.Errorf("could not resolve client info: %v", err)
				return
			}

			clientID = client.GetID()
		}

		return
	}

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		return errors.NewResponse(err, 500)
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		msg := re.Description
		if msg == "" {
			msg = re.Error.Error()
		}

		logger.Default().Warn(msg)
	})

	return srv
}
