package external

import (
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
)

// We're expecting that our users will be able to complete
// external auth loop in 15 minutes.
const (
	gothMaxSessionStoreAge = 60 * 15 // In seconds.
)

func setupGoth(eas *externalAuthSettings) {
	if eas == nil || !eas.enabled {
		log().Info("external authentication disabled")
		return
	}

	store := sessions.NewCookieStore([]byte(eas.sessionStoreSecret))
	store.MaxAge(gothMaxSessionStoreAge)
	store.Options.Path = "/auth/external"
	store.Options.HttpOnly = true
	store.Options.Secure = eas.sessionStoreSecure
	gothic.Store = store

	setupGothProviders(eas)

}

func setupGothProviders(eas *externalAuthSettings) {
	var (
		err    error
		scopes = []string{"email"}
	)

	// Purge all previously configured providers
	if l := len(goth.GetProviders()); l > 0 {
		log().Debug("removing existing providers", zap.Int("count", l))
		goth.ClearProviders()
	}

	var enabled = 0
	for _, pc := range eas.providers {
		if pc.enabled {
			enabled++
		}
	}

	log().Debug("initializing enabled external authentication providers", zap.Int("count", enabled))

	for name, pc := range eas.providers {
		var provider goth.Provider

		log := log().With(zap.String("provider", name))

		if !pc.enabled {
			continue
		}

		if strings.Index(name, "openid-connect.") == 0 {
			if pc.issuerUrl == "" {
				log.Error("failed to discover OIDC provider, URL empty")
				continue
			}

			wellKnown := strings.TrimSuffix(pc.issuerUrl, "/") + "/.well-known/openid-configuration"

			if provider, err = openidConnect.New(pc.key, pc.secret, pc.redirectUrl, wellKnown, scopes...); err != nil {
				log.Error("failed to discover OIDC provider", zap.Error(err), zap.String("well-known", wellKnown))
				continue
			} else {
				provider.SetName(name)
			}
		} else {
			switch name {
			case "github":
				provider = github.New(pc.key, pc.secret, pc.redirectUrl, scopes...)
			case "facebook":
				provider = facebook.New(pc.key, pc.secret, pc.redirectUrl, scopes...)
			case "gplus":
				provider = gplus.New(pc.key, pc.secret, pc.redirectUrl, scopes...)
			case "linkedin":
				provider = linkedin.New(pc.key, pc.secret, pc.redirectUrl, scopes...)
			}
		}

		if provider != nil {
			log.Info("external authentication provider added")
			goth.UseProviders(provider)
		}
	}
}
