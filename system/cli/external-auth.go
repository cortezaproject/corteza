package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/crusttech/crust/internal/settings"
	"github.com/crusttech/crust/system/internal/auth/external"
)

// Will perform OpenID connect auto-configuration
func externalAuthCmd(ctx context.Context, settingsService settings.Service) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "external-auth",
		Short: "External authentication",
	}

	autoDiscoverCmd := &cobra.Command{
		Use:   "auto-discovery [name] [url]",
		Short: "Auto discovers new OIDC client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var name, url = args[0], args[1]

			if eas, err := external.ExternalAuthSettings(settingsService); err != nil {
				exit(cmd, err)
			} else if eap, err := external.RegisterNewOpenIdClient(ctx, eas, name, url); err != nil {
				exit(cmd, err)
			} else if vv, err := eap.MakeValueSet("openid-connect." + name); err != nil {
				exit(cmd, err)
			} else if err := settingsService.BulkSet(vv); err != nil {
				exit(cmd, err)
			}
		},
	}

	cmd.AddCommand(autoDiscoverCmd)

	return cmd
}
