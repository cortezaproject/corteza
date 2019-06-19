package commands

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/internal/rand"
	"github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/internal/service"
)

func Settings(ctx context.Context) *cobra.Command {
	var (
		systemApiUrl, authFrontendUrl, authFromAddress, authFromName string

		cmd = &cobra.Command{
			Use:   "settings",
			Short: "Settings management",
		}
	)

	auto := &cobra.Command{
		Use:   "auto-configure",
		Short: "Run autoconfiguration",
		Run: func(cmd *cobra.Command, args []string) {
			_, _ = service.DefaultSettings.LoadAuthSettings()

			SettingsAutoConfigure(
				cmd,
				systemApiUrl,
				authFrontendUrl,
				authFromAddress,
				authFromName,
			)
		},
	}

	auto.Flags().StringVar(
		&systemApiUrl,
		"system-api-url",
		"",
		"System API URL (http://sytem.api.example.tld)",
	)

	auto.Flags().StringVar(
		&authFrontendUrl,
		"auth-frontend-url",
		"",
		"http://example.tld",
	)

	auto.Flags().StringVar(
		&authFromAddress,
		"auth-from-address",
		"",
		"name@example.tld",
	)

	auto.Flags().StringVar(
		&authFromName,
		"auth-from-name",
		"",
		"Name Surname",
	)

	list := &cobra.Command{
		Use:   "list",
		Short: "List all",
		Run: func(cmd *cobra.Command, args []string) {
			prefix := cmd.Flags().Lookup("prefix").Value.String()
			if kv, err := service.DefaultIntSettings.FindByPrefix(prefix); err != nil {
				cli.HandleError(err)
			} else {
				var maxlen int
				for _, v := range kv {
					if l := len(v.Name); l > maxlen {
						maxlen = l
					}
				}

				for _, v := range kv {
					cmd.Printf("%s%s\t%v\n", v.Name, strings.Repeat(" ", maxlen-len(v.Name)), v.Value)
				}
			}
		},
	}

	list.Flags().String("prefix", "", "Filter settings by prefix")

	get := &cobra.Command{
		Use: "get [key to get, ...]",

		Short: "Get value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if v, err := service.DefaultIntSettings.Get(args[0], 0); err != nil {
				cli.HandleError(err)
			} else if v != nil {
				cmd.Printf("%v\n", v.Value)
			}
		},
	}

	set := &cobra.Command{
		Use:   "set [key to set] [value]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			value := args[1]
			v := &settings.Value{
				Name: args[0],
			}

			if err := v.SetValueAsString(value); err != nil {
				cli.HandleError(err)
			}

			cli.HandleError(service.DefaultIntSettings.Set(v))
		},
	}

	imp := &cobra.Command{
		Use:   "import [file]",
		Short: "Import settings as JSON from stdin or file",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				fh  *os.File
				err error
			)

			if len(args) > 0 {
				fh, err = os.Open(args[0])
				cli.HandleError(err)
			} else {
				fh = os.Stdin
			}

			var (
				decoder = json.NewDecoder(fh)
				input   = map[string]interface{}{}
				vv      settings.ValueSet
			)

			cli.HandleError(decoder.Decode(&input))

			for k, v := range input {
				val := &settings.Value{Name: k}

				cli.HandleError(val.SetValue(v))
				vv = append(vv, val)
			}

			if len(vv) > 0 {
				cli.HandleError(service.DefaultIntSettings.BulkSet(vv))
			}
		},
	}

	exp := &cobra.Command{
		Use:   "export [file]",
		Short: "Import settings as JSON to stdout or file",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				fh  *os.File
				err error
			)

			if len(args) > 0 {
				fh, err = os.Create(args[0])
				cli.HandleError(err)
			} else {
				fh = os.Stdin
			}

			var (
				encoder = json.NewEncoder(fh)
			)

			encoder.SetIndent("", "  ")

			if vv, err := service.DefaultIntSettings.FindByPrefix(""); err != nil {
				cli.HandleError(err)
			} else {
				cli.HandleError(encoder.Encode(vv.KV()))
			}
		},
	}

	del := &cobra.Command{
		Use:   "delete [keys, ...]",
		Short: "Set value (raw JSON) for a specific key",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for a := 0; a < len(args); a++ {
				cli.HandleError(service.DefaultIntSettings.Delete(args[a], 0))
			}
		},
	}

	cmd.AddCommand(
		auto,
		list,
		get,
		set,
		del,
		imp,
		exp,
	)

	return cmd
}

func SettingsAutoConfigure(cmd *cobra.Command, systemApiUrl, frontendUrl, fromAddress, fromName string) {
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
				cmd.Printf("could not marshal setting value: %v", err)
				return
			}
		}

		err = service.DefaultIntSettings.Set(v)
		if err != nil {
			cmd.Printf("could not store setting: %v", err)
		}
	}

	setIfMissing := func(name string, value interface{}) {
		if existing, err := service.DefaultIntSettings.Get(name, 0); err == nil && existing == nil {
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
		extRedirUrl, _ := service.DefaultIntSettings.GetGlobalString("auth.external.redirect-url")
		return strings.Contains(extRedirUrl, "https://")
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

		setIfMissing("auth.frontend.url.base", func() interface{} {
			return frontendUrl + "/"
		})
	}

	// Auth email (password reset, email confirmation)
	setIfMissing("auth.mail.from-address", func() interface{} {
		if len(fromAddress) > 0 {
			return fromAddress
		}
		return "change-me@example.tld"
	})

	setIfMissing("auth.mail.from-name", func() interface{} {
		if len(fromName) > 0 {
			return fromName
		}

		return "Corteza Team"
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
