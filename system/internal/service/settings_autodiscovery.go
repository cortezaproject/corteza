package service

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/rand"
	internalSettings "github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/pkg/api"
)

// Discovers "auth.%" settings from the environment
//
// when other kinds of auto-discoverable settings come, lambdas inside will probably need a bit of refactoring
func authSettingsAutoDiscovery(log *zap.Logger, current internalSettings.ValueSet) (internalSettings.ValueSet, error) {
	type (
		stringWrapper func() string
		boolWrapper   func() bool
	)

	if log == nil {
		log = zap.NewNop()
	}

	log = log.Named("discovery")

	var (
		new = current

		// Setter
		//
		// Finds existing settings, tries with environmental "PROVISION_SETTINGS_AUTH_..." probing
		// and falls back to default value
		//
		// We are extremely verbose here - we want to show all the info available and
		// how settings were discovered and set
		set = func(name string, env string, def interface{}) {
			var (
				log = log.With(
					zap.String("name", name),
				)

				v     = current.First(name)
				value interface{}
			)

			if v != nil {
				// Nothing to discover, already set
				log.Info("already set", zap.Any("value", v.String()))
				return
			}

			v = &internalSettings.Value{Name: name}

			value, envExists := os.LookupEnv(env)

			switch dfn := def.(type) {
			case stringWrapper:
				log = log.With(zap.String("type", "string"))
				// already a string, no need to do any magic
				if envExists {
					log = log.With(zap.String("env", env), zap.Any("value", value))
				} else {
					value = dfn()
					log = log.With(zap.Any("default", value))
				}
			case boolWrapper:
				log = log.With(zap.String("type", "bool"))

				if envExists {
					value = cast.ToBool(value)
					log = log.With(zap.String("env", env), zap.Any("value", value))
				} else {
					value = dfn()
					log = log.With(zap.Any("default", value))
				}

			default:
				log.Error("unsupported type")
				return
			}

			if err := v.SetValue(value); err != nil {
				log.Error("could not set value", zap.Error(err))
				return
			}

			log.Info("value auto-discovered")

			new.Replace(v)
		}

		// Default value functions
		//
		// all are wrapped (stringWrapper, boolWrapper) to delay execution
		// of the function to the very last point

		frontendUrl = func(path string) stringWrapper {
			const (
				feBase   = "auth.frontend.url.base"
				extRedir = "auth.external.redirect-url"
			)

			return func() (base string) {
				base = new.First(feBase).String()

				if len(base) == 0 {

					// Not found, try to get it from the external redirect URL
					redirURL := new.First(extRedir).String()
					if len(redirURL) == 0 {
						return
					}

					log.Info(
						"discovering frontend url from '"+extRedir+"'",
						zap.String(extRedir, redirURL))

					// Removing placeholder
					redirURL = fmt.Sprintf(redirURL, "")

					p, err := url.Parse(redirURL)
					if err != nil {
						log.Error("could not parse '"+extRedir+"'", zap.Error(err))
						return
					}

					h := p.Host
					s := "api."
					if i := strings.Index(h, s); i > 0 {
						// If there is a "api." prefix in the hostname of the external redirect-uri value
						// cut it off and use that as a frontend url base
						h = h[i+len(s):]
					}

					base = p.Scheme + "://" + h
				}

				if len(base) > 0 {
					return strings.TrimRight(base, "/") + path
				}

				return ""
			}
		}

		// Assuming secure backend when redirect URL starts with https://
		isSecure = func() boolWrapper {
			return func() bool {
				return strings.Index(new.First("auth.external.redirect-url").String(), "https://") == 0
			}
		}

		// Assume we have emailing capabilities if SMTP_HOST variable is set
		emailCapabilities = func() boolWrapper {
			return func() bool {
				val, has := os.LookupEnv("SMTP_HOST")
				return has && len(val) > 0
			}
		}

		// Where should external authentication providers redirect to?
		// we need to set full, absolute URL to the callback endpoint
		externalAuthRedirectUrl = func() stringWrapper {
			return func() string {
				var (
					path = "/auth/external/%s/callback"

					// All env keys we'll check, first that has any value set, will be used as hostname
					keysWithHostnames = []string{
						"DOMAIN",
						"LETSENCRYPT_HOST",
						"VIRTUAL_HOST",
						"HOSTNAME",
						"HOST",
					}
				)

				// Prefix path if we're running wrapped as a monolith:
				if api.Monolith {
					path = "/system" + path
				}

				// Finally, add any prefix
				path = strings.TrimRight(api.BaseURL, "/") + path

				log.Info("scanning env variables for hostname", zap.Strings("candidates", keysWithHostnames))

				for _, key := range keysWithHostnames {
					if host, has := os.LookupEnv(key); has {
						log.Info("hostname env variable found", zap.String("env", key))
						// Make life easier for development in local environment,
						// and set HTTP schema. Might cause problems if someone
						// is using valid external hostname
						if strings.Contains(host, "local.") {
							return "http://" + host + path
						} else {
							return "https://" + host + path
						}
					} else {
					}
				}

				// Fallback is empty string
				// this will cause error when doing OIDC auto-discovery (and we want that)
				// @todo ^^
				return ""
			}
		}

		rand stringWrapper = func() string {
			return string(rand.Bytes(64))
		}

		wrapBool = func(val bool) boolWrapper {
			return func() bool { return val }
		}

		wrapString = func(val string) stringWrapper {
			return func() string { return val }
		}
	)

	// List of name-value pairs we need to iterate and set
	list := []struct {
		// Setting name
		nme string

		// provision environmental variable name
		// we're using full variable name here so developers
		// can find where things are comming from
		env string

		// default value
		// expects one of the *wrapper() functions
		// this also determinate the value type of the setting and casting rules for the env value
		def interface{}
	}{
		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //
		// External auth

		// Enable external auth
		{
			"auth.external.enabled",
			"PROVISION_SETTINGS_AUTH_EXTERNAL_ENABLED",
			wrapBool(true)},

		{
			"auth.external.redirect-url",
			"PROVISION_SETTINGS_AUTH_EXTERNAL_REDIRECT_URL",
			externalAuthRedirectUrl()},

		{
			"auth.external.session-store-secret",
			"PROVISION_SETTINGS_AUTH_EXTERNAL_SESSION_STORE_SECRET",
			rand},

		// Disable external auth
		{
			"auth.external.session-store-secure",
			"PROVISION_SETTINGS_AUTH_EXTERNAL_SESSION_STORE_SECURE",
			isSecure()},

		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //
		// Auth frontend

		{
			"auth.frontend.url.base",
			"PROVISION_SETTINGS_AUTH_FRONTEND_URL_BASE",
			frontendUrl("/")},

		// @todo w/o token=
		{
			"auth.frontend.url.password-reset",
			"PROVISION_SETTINGS_AUTH_FRONTEND_URL_PASSWORD_RESET",
			frontendUrl("/auth/reset-password?token=")},

		// @todo w/o token=
		{
			"auth.frontend.url.email-confirmation",
			"PROVISION_SETTINGS_AUTH_FRONTEND_URL_EMAIL_CONFIRMATION",
			frontendUrl("/auth/confirm-email?token=")},

		// @todo check if this is correct?!
		{
			"auth.frontend.url.redirect",
			"PROVISION_SETTINGS_AUTH_FRONTEND_URL_REDIRECT",
			frontendUrl("/auth")},

		// Auth email
		{
			"auth.mail.from-address",
			"PROVISION_SETTINGS_AUTH_EMAIL_FROM_ADDRESS",
			wrapString("to-be-configured@example.tld")},

		{
			"auth.mail.from-name",
			"PROVISION_SETTINGS_AUTH_EMAIL_FROM_NAME",
			wrapString("Corteza Team (to-be-configured)")},

		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //
		// Enable internal login
		{
			"auth.internal.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_ENABLED",
			wrapBool(true)},

		// Enable internal signup
		{
			"auth.internal.signup.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_SIGNUP_ENABLED",
			wrapBool(true)},

		// Enable email confirmation if we have email capabilities
		{
			"auth.internal.signup-email-confirmation-required",
			"PROVISION_SETTINGS_AUTH_INTERNAL_SIGNUP_EMAIL_CONFIRMATION_REQUIRED",
			emailCapabilities()},

		// Enable password reset if we have email capabilities
		{
			"auth.internal.password-reset.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_PASSWORD_RESET_ENABLED",
			emailCapabilities()},
	}

	for _, item := range list {
		set(item.nme, item.env, item.def)
	}

	// return new, nil
	return current.Changed(new), nil
}
