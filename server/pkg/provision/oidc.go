package provision

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/auth/external"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
)

// Provisions OIDC providers from PROVISION_OIDC_PROVIDER env variable
//
// Env variable should contains space delimited pairs of providers (<name> <provider> ....)
func oidcAutoDiscovery(ctx context.Context, log *zap.Logger, s store.Settings, opt options.AuthOpt) (err error) {
	var provider = strings.TrimSpace(options.EnvString("PROVISION_OIDC_PROVIDER", ""))

	log.Debug("OIDC auto discovery provision",
		zap.String("envkey", "PROVISION_OIDC_PROVIDER"),
		zap.String("providers", provider),
	)

	if len(provider) == 0 {
		return
	}

	var (
		providers  = strings.Split(provider, " ")
		plen       = len(providers)
		name, purl string
		eap        *types.ExternalAuthProvider
	)

	if plen%2 == 1 {
		return fmt.Errorf("expecting even number of providers")
	}

	for p := 0; p < plen; p = p + 2 {
		name, purl = providers[p], providers[p+1]

		// force:    false
		// because we do not want to override the provider every time the system restarts
		//
		// validate: false
		// because at the initial (empty db) start, we cannot actually validate (server is not yet up)
		//
		// enable:   true
		// we want provider & the entire external auth to be validated
		eap, err = external.RegisterOidcProvider(ctx, log, s, opt, name, purl, false, false, true)

		if err != nil {
			log.Error(
				"could not register OIDC provider",
				zap.String("url", purl),
				zap.String("name", name),
				zap.Error(err))
			return
		} else if eap == nil {
			log.Info("provider already exists",
				zap.String("name", name))
		} else {
			log.Info("provider successfully registered",
				zap.String("url", purl),
				zap.String("key", eap.Key),
				zap.String("name", name))
		}
	}

	return
}

func authAddExternals(ctx context.Context, log *zap.Logger, s store.Settings) (err error) {
	var (
		kinds = []string{
			"github",
			"facebook",
			"google",
			"linkedin",
			"oidc",
		}

		key, name string

		pp []string

		eap *types.ExternalAuthProvider
	)

	for _, kind := range kinds {
		key = "PROVISION_SETTINGS_AUTH_EXTERNAL_" + strings.ToUpper(kind)
		val, has := os.LookupEnv(key)

		if !has {
			continue
		}

		if val = strings.TrimSpace(val); val == "" {
			log.Warn(fmt.Sprintf("found empty value for %s", key))
			continue
		} else {
			log.Info(fmt.Sprintf("provosioning external auth provider from %s", key))
		}

		eap = &types.ExternalAuthProvider{Enabled: true}

		if kind == "oidc" {
			pp = strings.SplitN(val, " ", 4)

			if len(pp) < 4 {
				log.Warn(fmt.Sprintf("expecting 4 values (name, issuer-url, ket, secret) in %s", key))
			}

			// Spread name, issuer-url, key and secret from provision string for OIDC provider
			name, eap.IssuerUrl, eap.Key, eap.Secret = pp[0], pp[1], pp[2], pp[3]

			eap.Handle = external.OIDC_PROVIDER_PREFIX + name
		} else {
			pp = strings.SplitN(val, " ", 2)

			// Spread key and secret from provision string
			eap.Key, eap.Secret = pp[0], pp[1]
			eap.Handle = kind
		}

		_ = external.AddProvider(ctx, log, s, eap, false)
	}

	return
}
