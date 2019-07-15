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

	"github.com/cortezaproject/corteza-server/system/internal/service"
)

// We're expecting that our users will be able to complete
// external auth loop in 15 minutes.
const (
	gothMaxSessionStoreAge = 60 * 15 // In seconds.

	WellKnown = "/.well-known/openid-configuration"
)

func setupGoth(as service.AuthSettings) {
	if !as.ExternalEnabled {
		log().Info("external authentication disabled")
		return
	}

	store := sessions.NewCookieStore([]byte(as.ExternalSessionStoreSecret))
	store.MaxAge(gothMaxSessionStoreAge)
	store.Options.HttpOnly = true
	store.Options.Secure = as.ExternalSessionStoreSecure
	gothic.Store = store

	log().Debug("registering cookie session store")

	if store.Options.Secure {
		log().Debug("cookie session store has 'secure' flag ON, make sure this URL is accessed via HTTPS")

	}

	setupGothProviders(as)

}

func setupGothProviders(as service.AuthSettings) {
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
	for _, pc := range as.ExternalProviders {
		if pc.Enabled {
			enabled++
		}
	}

	log().Debug("initializing enabled external authentication providers", zap.Int("count", enabled))

	for name, pc := range as.ExternalProviders {
		var provider goth.Provider

		log := log().With(zap.String("provider", name))

		if !pc.Enabled {
			continue
		}

		if strings.Index(name, "openid-connect.") == 0 {
			if pc.IssuerUrl == "" {
				log.Error("failed to discover OIDC provider, URL empty")
				continue
			}

			wellKnown := strings.TrimSuffix(pc.IssuerUrl, "/") + WellKnown

			if provider, err = openidConnect.New(pc.Key, pc.Secret, pc.RedirectUrl, wellKnown, scopes...); err != nil {
				log.Error("failed to discover OIDC provider", zap.Error(err), zap.String("well-known", wellKnown))
				continue
			} else {
				provider.SetName(name)
			}
		} else {
			switch name {
			case "github":
				provider = github.New(pc.Key, pc.Secret, pc.RedirectUrl, scopes...)
			case "facebook":
				provider = facebook.New(pc.Key, pc.Secret, pc.RedirectUrl, scopes...)
			case "gplus":
				provider = gplus.New(pc.Key, pc.Secret, pc.RedirectUrl, scopes...)
			case "linkedin":
				provider = linkedin.New(pc.Key, pc.Secret, pc.RedirectUrl, scopes...)
			}
		}

		if provider != nil {
			log.Info("external authentication provider added")
			goth.UseProviders(provider)
		}
	}
}
