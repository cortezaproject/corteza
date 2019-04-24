package cli

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/rand"
	"github.com/crusttech/crust/internal/settings"
	systemService "github.com/crusttech/crust/system/internal/service"
)

func settingsCmd(ctx context.Context, setSvc settings.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "Settings management",
	}

	auto := &cobra.Command{
		Use:   "auto-configure",
		Short: "Run autoconfiguration",
		Run: func(cmd *cobra.Command, args []string) {
			systemService.DefaultSettings.LoadAuthSettings()

			settingsAutoConfigure(
				setSvc,
				cmd.Flags().Lookup("system-api-url").Value.String(),
				cmd.Flags().Lookup("auth-frontend-url").Value.String(),
				cmd.Flags().Lookup("auth-from-address").Value.String(),
				cmd.Flags().Lookup("auth-from-address").Value.String(),
			)
		},
	}

	auto.Flags().String("system-api-url", "", "System API URL (http://sytem.api.example.tld)")
	auto.Flags().String("auth-frontend-url", "", "http://example.tld")
	auto.Flags().String("auth-from-address", "", "name@example.tld")
	auto.Flags().String("auth-from-name", "", "Name Surname")

	list := &cobra.Command{
		Use:   "list",
		Short: "List all",
		Run: func(cmd *cobra.Command, args []string) {
			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := setSvc.FindByPrefix(prefix); err != nil {
				exit(cmd, err)
			} else {
				for _, v := range kv {
					cmd.Printf("%s\t%v\n", v.Name, v.Value)
				}
			}
		},
	}

	list.Flags().String("prefix", "", "Filter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get]",

		Short: "Get value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := setSvc.Get(args[0], 0); err != nil {
				exit(cmd, err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
			exit(cmd, nil)
		},
	}

	set := &cobra.Command{
		Use:   "set [key to set] [value",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			value := args[1]
			v := &settings.Value{
				Name: args[0],
			}

			if err := v.SetValueAsString(value); err != nil {
				exit(cmd, err)
			}

			exit(cmd, setSvc.Set(v))
		},
	}

	del := &cobra.Command{
		Use:   "delete [key to remove]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			exit(cmd, setSvc.Delete(args[0], 0))
		},
	}

	cmd.AddCommand(
		auto,
		list,
		get,
		set,
		del,
	)

	return cmd
}

func settingsAutoConfigure(setSvc settings.Service, systemApiUrl, frontendUrl, fromAddress, fromName string) {
	set := func(name string, value interface{}) {
		var (
			v   *settings.Value
			ok  bool
			err error
		)

		if setFn, ok := value.(func() interface{}); ok {
			// No existing value found
			if value = setFn(); value == nil {
				// setFn returned nil, skip everything..
				return
			}
		}

		// Did we get something else than *settings.Value?
		// if so, wrap/marshal it into *settings.Value
		if v, ok = value.(*settings.Value); !ok {
			v = &settings.Value{Name: name}
			if v.Value, err = json.Marshal(value); err != nil {
				log.Printf("could not marshal setting value: %v", err)
				return
			}
		}

		err = setSvc.Set(v)
		if err != nil {
			log.Printf("could not store setting: %v", err)
		}
	}

	setIfMissing := func(name string, value interface{}) {
		if existing, err := setSvc.Get(name, 0); err == nil && existing == nil {
			set(name, value)
		}
	}

	// Where should external authentication providers redirect to?
	setIfMissing("auth.external.redirect-url", func() interface{} {
		path := "/auth/external/%s/callback"

		if len(systemApiUrl) > 0 {
			if strings.Index(systemApiUrl, "http") != 0 {
				return "https://" + systemApiUrl + path
			} else {
				return systemApiUrl + path
			}
		}

		if leHost, has := os.LookupEnv("LETSENCRYPT_HOST"); has {
			return "https://" + leHost + path
		} else if vHost, has := os.LookupEnv("VIRTUAL_HOST"); has {
			return "http://" + vHost + path
		} else {
			// Fallback to local
			return "http://system.api.local.crust.tech" + path
		}
	})

	setIfMissing("auth.external.session-store-secret", func() interface{} {
		// Generate session store secret if missing
		return string(rand.Bytes(64))
	})

	setIfMissing("auth.external.session-store-secure", func() interface{} {
		// Try to determines if we need secure session store from redirect URL scheme
		extRedirUrl, _ := setSvc.GetGlobalString("auth.external.redirect-url")
		return strings.Index(extRedirUrl, "https://") > -1
	})

	if len(frontendUrl) > 0 {
		setIfMissing("auth.frontend.url.password-reset", func() interface{} {
			return frontendUrl + "/auth/reset-password?token="
		})

		setIfMissing("auth.frontend.url.email-confirmation", func() interface{} {
			return frontendUrl + "/auth/confirm-email?token="
		})

		setIfMissing("auth.frontend.url.redirect", func() interface{} {
			return frontendUrl + "/auth/"
		})
	}

	// Auth email (password reset, email confirmation)
	setIfMissing("auth.mail.from-address", func() interface{} {
		if len(fromAddress) > 0 {
			return fromAddress
		}
		return "change-me@local.crust.tech"
	})

	setIfMissing("auth.mail.from-name", func() interface{} {
		if len(fromName) > 0 {
			return fromName
		}

		return "Crust"
	})

	// No external providers preconfigured, so disable
	setIfMissing("auth.external.enabled", false)

	// Enable internal by default
	setIfMissing("auth.internal.enabled", true)

	// Enable user creation
	setIfMissing("auth.internal.signup.enabled", true)

	// No need to confirm email
	setIfMissing("auth.internal.signup-email-confirmation-required", false)

	// We need password reset
	setIfMissing("auth.internal.password-reset.enabled", true)
}
