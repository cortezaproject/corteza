package external

import (
	"strings"

	"github.com/cortezaproject/corteza/server/auth/settings"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/openidConnect"
	"go.uber.org/zap"
)

// We're expecting that our users will be able to complete
// external auth loop in 15 minutes.
const (
	WellKnown = "/.well-known/openid-configuration"
)

func SetupGothProviders(log *zap.Logger, redirectUrl string, ep ...settings.Provider) {
	var (
		err error
	)

	// Purge all previously configured providers
	if l := len(goth.GetProviders()); l > 0 {
		log.Debug("removing existing providers", zap.Int("count", l))
		goth.ClearProviders()
	}

	log.Debug("initializing enabled external authentication providers", zap.Int("count", len(ep)))

	for _, pc := range ep {
		var provider goth.Provider

		log := log.With(zap.String("provider", pc.Handle))

		redirect := pc.RedirectUrl
		if redirect == "" {
			// If redirect URL is not explicitly set for this provider,
			// generate one from template string
			redirect = strings.Replace(redirectUrl, "{provider}", pc.Handle, 1)
		}

		if strings.HasPrefix(pc.Handle, OIDC_PROVIDER_PREFIX) {
			if pc.IssuerUrl == "" {
				log.Error("failed to discover OIDC provider, URL empty")
				continue
			}

			wellKnown := strings.TrimSuffix(pc.IssuerUrl, "/") + WellKnown

			var scope []string
			if len(pc.Scope) > 0 {
				scope = strings.Split(pc.Scope, " ")
			} else {
				scope = append(scope, "email")
			}

			if provider, err = openidConnect.New(pc.Key, pc.Secret, redirect, wellKnown, scope...); err != nil {
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
			log.Info("external authentication provider added", zap.String("key", pc.Key))
			goth.UseProviders(provider)
		}
	}
}
