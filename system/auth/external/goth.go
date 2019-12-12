package external

import (
	"fmt"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/system/types"
)

// We're expecting that our users will be able to complete
// external auth loop in 15 minutes.
const (
	gothMaxSessionStoreAge = 60 * 15 // In seconds.

	WellKnown = "/.well-known/openid-configuration"
)

func setupGoth(s *types.Settings) {
	if !s.Auth.External.Enabled {
		log().Info("external authentication disabled")
		return
	}

	store := sessions.NewCookieStore([]byte(s.Auth.External.SessionStoreSecret))
	store.MaxAge(gothMaxSessionStoreAge)
	store.Options.HttpOnly = true
	store.Options.Secure = s.Auth.External.SessionStoreSecure
	gothic.Store = store

	log().Debug("registering cookie session store")

	if store.Options.Secure {
		log().Debug("cookie session store has 'secure' flag ON, make sure this URL is accessed via HTTPS")

	}

	setupGothProviders(s)

}

func setupGothProviders(s *types.Settings) {
	var (
		err error
	)

	// Purge all previously configured providers
	if l := len(goth.GetProviders()); l > 0 {
		log().Debug("removing existing providers", zap.Int("count", l))
		goth.ClearProviders()
	}

	var enabled = 0
	for _, pc := range s.Auth.External.Providers {
		if pc.Enabled {
			enabled++
		}
	}

	log().Debug("initializing enabled external authentication providers", zap.Int("count", enabled))

	for _, pc := range s.Auth.External.Providers {
		var provider goth.Provider

		log := log().With(zap.String("provider", pc.Handle))

		if !pc.Enabled {
			continue
		}

		redirect := pc.RedirectUrl
		if redirect == "" {
			// If redirect URL is not explicitly set for this provider,
			// generate one from template string
			redirect = fmt.Sprintf(s.Auth.External.RedirectUrl, pc.Handle)
		}

		if strings.Index(pc.Handle, OIDC_PROVIDER_PREFIX) == 0 {
			if pc.IssuerUrl == "" {
				log.Error("failed to discover OIDC provider, URL empty")
				continue
			}

			wellKnown := strings.TrimSuffix(pc.IssuerUrl, "/") + WellKnown

			if provider, err = openidConnect.New(pc.Key, pc.Secret, redirect, wellKnown, "email"); err != nil {
				log.Error("failed to discover OIDC provider", zap.Error(err), zap.String("well-known", wellKnown))
				continue
			} else {
				provider.SetName(pc.Handle)
			}
		} else {
			switch pc.Handle {
			case "github":
				provider = github.New(pc.Key, pc.Secret, redirect, "user:email")
			case "facebook":
				provider = facebook.New(pc.Key, pc.Secret, redirect, "email")
			case "google":
				provider = google.New(pc.Key, pc.Secret, redirect, "email")
			case "linkedin":
				provider = linkedin.New(pc.Key, pc.Secret, redirect, "email")
			}
		}

		if provider != nil {
			log.Info(
				"external authentication provider added",
				zap.String("key", pc.Key),
				zap.String("redirectUrl", redirect),
			)
			goth.UseProviders(provider)
		}
	}
}
