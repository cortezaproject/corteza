package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/auth/external"
)

// Will perform OpenID connect auto-configuration
func ExternalAuth(ctx context.Context, rootCmd *cobra.Command, settingsService settings.Service) {
	exit := func(err error) {
		if err != nil {
			rootCmd.Printf("Error: %v\n", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}

	autoDiscover := &cobra.Command{
		Use:   "auto-discovery [name] [url]",
		Short: "Auto discovers new OIDC client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var name, url = args[0], args[1]

			if eas, err := external.ExternalAuthSettings(settingsService); err != nil {
				exit(err)
			} else if eap, err := external.RegisterNewOpenIdClient(ctx, eas, name, url); err != nil {
				exit(err)
			} else if vv, err := eap.MakeValueSet("openid-connect." + name); err != nil {
				exit(err)
			} else if err := settingsService.BulkSet(vv); err != nil {
				exit(err)
			}
		},
	}

	settingsCmd := &cobra.Command{
		Use:   "external-auth",
		Short: "External authentication",
	}

	settingsCmd.AddCommand(autoDiscover)

	rootCmd.AddCommand(settingsCmd)
}
