package external

import (
	"log"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/openidConnect"
)

// We're expecting that our users will be able to complete
// external auth loop in 15 minutes.
const (
	gothMaxSessionStoreAge = 60 * 15 // In seconds.
)

func setupGoth(eas *externalAuthSettings) {
	if eas == nil || !eas.enabled {
		log.Printf("external authentication disabled (%v)", eas)
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
		log.Println("removing existing providers", l)
		goth.ClearProviders()
	}

	var enabled uint
	for _, pc := range eas.providers {
		if pc.enabled {
			enabled++
		}
	}

	log.Printf("initializing external authentication providers (%d)", enabled)

	for name, pc := range eas.providers {
		var provider goth.Provider

		if !pc.enabled {
			continue
		}

		if strings.Index(name, "openid-connect.") == 0 {
			if pc.issuerUrl == "" {
				log.Printf("failed to discover %q (issuer URL is empty)", name)
				continue
			}

			wellKnown := strings.TrimSuffix(pc.issuerUrl, "/") + "/.well-known/openid-configuration"

			if provider, err = openidConnect.New(pc.key, pc.secret, pc.redirectUrl, wellKnown, scopes...); err != nil {
				log.Printf("failed to discover %q (auto discovery URL: %s) OIDC provider: %v", name, wellKnown, err)
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
			log.Printf("external authentication provider %q added", name)
			goth.UseProviders(provider)
		}
	}
}
