package commands

import (
	"github.com/cortezaproject/corteza-server/auth/external"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/spf13/cobra"
)

// Will perform OpenID connect auto-configuration
func Auth(app serviceInitializer, opt options.AuthOpt) *cobra.Command {
	var (
		enableDiscoveredProvider               bool
		skipValidationOnAutoDiscoveredProvider bool
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "External authentication",
	}

	autoDiscoverCmd := &cobra.Command{
		Use:     "auto-discovery [name] [url]",
		Short:   "Auto discovers new OIDC client",
		Args:    cobra.ExactArgs(2),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := auth.SetSuperUserContext(cli.Context())
			_, err := external.RegisterOidcProvider(
				ctx,
				opt,
				args[0],
				args[1],
				true,
				!skipValidationOnAutoDiscoveredProvider,
				enableDiscoveredProvider,
			)

			cli.HandleError(err)

			if enableDiscoveredProvider {
				cmd.Println("OIDC provider successfully added and enabled.")
			} else {
				cmd.Println("OIDC provider successfully added (still disabled).")
			}
		},
	}

	autoDiscoverCmd.Flags().BoolVar(
		&enableDiscoveredProvider,
		"enable",
		false,
		"Enable this provider and external auth")

	autoDiscoverCmd.Flags().BoolVar(
		&skipValidationOnAutoDiscoveredProvider,
		"skip-validation",
		false,
		"Skip validation")

	jwtCmd := &cobra.Command{
		Use:     "jwt [email-or-id]",
		Short:   "Generates new JWT for a user",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())

				user *types.User
				err  error

				userStr = args[0]
			)

			user, err = service.DefaultUser.FindByAny(ctx, userStr)
			cli.HandleError(err)

			err = service.DefaultAuth.LoadRoleMemberships(ctx, user)
			cli.HandleError(err)

			cmd.Println(auth.DefaultJwtHandler.Encode(user))
		},
	}

	testEmails := &cobra.Command{
		Use:     "test-notifications [recipient]",
		Short:   "Sends samples of all authentication notification to receipient",
		Args:    cobra.ExactArgs(1),
		PreRunE: commandPreRunInitService(app),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ctx = auth.SetSuperUserContext(cli.Context())
				err error
				ntf = service.DefaultAuthNotification
			)

			// Update current settings to be sure that we do not have outdated values
			cli.HandleError(service.DefaultSettings.UpdateCurrent(ctx))

			err = ntf.PasswordReset(ctx, "en", args[0], "notification-testing-token")
			cli.HandleError(err)
		},
	}

	cmd.AddCommand(
		autoDiscoverCmd,
		testEmails,
		jwtCmd,
	)

	return cmd
}
