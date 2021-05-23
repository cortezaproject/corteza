package provision

import (
	"context"
	"fmt"
	"os"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Discovers "auth.%" settings from the environment
//
// when other kinds of auto-discoverable settings come, lambdas inside will probably need a bit of refactoring
func authSettingsAutoDiscovery(ctx context.Context, log *zap.Logger, s store.Storer) (err error) {
	type (
		stringWrapper func() string
		boolWrapper   func() bool
	)

	log = log.Named("auth.settings-discovery")

	var (
		// Setter
		//
		// Finds existing settings, tries with environmental "PROVISION_SETTINGS_AUTH_..." probing
		// and falls back to default value
		//
		// We are extremely verbose here - we want to show all the info available and
		// how settings were discovered and set
		//
		// @todo generalize and move under settings
		set = func(name string, env string, def interface{}, maskSensitive bool) error {

			var (
				log = log.With(
					zap.String("name", name),
				)

				envExists bool
				value     interface{}

				v, err = s.LookupSettingByNameOwnedBy(ctx, name, 0)
			)

			if !errors.IsNotFound(err) && err != nil {
				return fmt.Errorf("could not load settings value for '%s': %w", name, err)
			}

			if v != nil {
				// Nothing to discover, already set
				log.Debug("already set", logger.MaskIf("value", v, maskSensitive))
				return nil
			}

			v = &types.SettingValue{Name: name}

			value, envExists = os.LookupEnv(env)

			switch dfn := def.(type) {
			case stringWrapper:
				log = log.With(zap.String("type", "string"))
				// already a string, no need to do any magic
				if envExists {
					log = log.With(zap.String("env", env), logger.MaskIf("value", value, maskSensitive))
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
				return fmt.Errorf("unsupported type %T for '%s'", def, name)
			}

			if err := v.SetValue(value); err != nil {
				return fmt.Errorf("could not set value to '%q': %w", name, err)
			}

			log.Debug("value auto-discovered")
			return s.UpsertSetting(ctx, v)
		}

		// Assume we have emailing capabilities if SMTP_HOST variable is set
		emailCapabilities = func() boolWrapper {
			return func() bool {
				val, has := os.LookupEnv("SMTP_HOST")
				return has && len(val) > 0
			}
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
		// can find where things are coming from
		env string

		// default value
		// expects one of the *wrapper() functions
		// this also determinate the value type of the setting and casting rules for the env value
		def interface{}

		// mask value if sensitive
		mask bool
	}{
		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //
		// External auth

		// Enable federated auth
		{
			"auth.external.enabled",
			"PROVISION_SETTINGS_AUTH_EXTERNAL_ENABLED",
			wrapBool(true),
			false},

		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //

		// Auth email
		{
			"auth.mail.from-address",
			"PROVISION_SETTINGS_AUTH_EMAIL_FROM_ADDRESS",
			wrapString("info@example.tld"),
			false},

		{
			"auth.mail.from-name",
			"PROVISION_SETTINGS_AUTH_EMAIL_FROM_NAME",
			wrapString("Example Sender"),
			false},

		// // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // // //
		// Enable internal login
		{
			"auth.internal.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_ENABLED",
			wrapBool(true),
			false},

		// Enable internal signup
		{
			"auth.internal.signup.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_SIGNUP_ENABLED",
			wrapBool(true),
			false},

		// Enable email confirmation if we have email capabilities
		{
			"auth.internal.signup.email-confirmation-required",
			"PROVISION_SETTINGS_AUTH_INTERNAL_SIGNUP_EMAIL_CONFIRMATION_REQUIRED",
			emailCapabilities(),
			false},

		// Enable password reset if we have email capabilities
		{
			"auth.internal.password-reset.enabled",
			"PROVISION_SETTINGS_AUTH_INTERNAL_PASSWORD_RESET_ENABLED",
			emailCapabilities(),
			false},
	}

	for _, item := range list {
		if err = set(item.nme, item.env, item.def, item.mask); err != nil {
			return err
		}
	}

	return nil
}
