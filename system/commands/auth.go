package commands

import (
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/auth/external"
	"github.com/cortezaproject/corteza-server/system/repository"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Will perform OpenID connect auto-configuration
func Auth() *cobra.Command {
	var (
		enableDiscoveredProvider               bool
		skipValidationOnAutoDiscoveredProvider bool
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "External authentication",
	}

	autoDiscoverCmd := &cobra.Command{
		Use:   "auto-discovery [name] [url]",
		Short: "Auto discovers new OIDC client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())
			)

			_, err := external.RegisterOidcProvider(
				ctx,
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
		Use:   "jwt [email-or-id]",
		Short: "Generates new JWT for a user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())
				db  = factory.Database.MustGet("system", "default")

				userRepo = repository.User(ctx, db)
				roleRepo = repository.Role(ctx, db)
				// authSvc  = service.Auth(ctx)

				user *types.User
				err  error
				ID   uint64
				rr   types.RoleSet

				userStr = args[0]
			)

			if user, err = userRepo.FindByEmail(userStr); repository.ErrUserNotFound.Eq(err) {
				if regexp.MustCompile(`/^\d+$/`).MatchString(userStr) {
					if ID, err = strconv.ParseUint(userStr, 10, 64); err == nil {
						user, err = userRepo.FindByID(ID)
					}
				}
			}

			if err == nil {
				rr, _, err = roleRepo.Find(types.RoleFilter{MemberID: user.ID})
			}

			if err != nil {
				cli.HandleError(err)
			}

			user.SetRoles(rr.IDs())

			cmd.Println(auth.DefaultJwtHandler.Encode(user))
		},
	}

	testEmails := &cobra.Command{
		Use:   "test-notifications [recipient]",
		Short: "Sends samples of all authentication notification to receipient",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ctx = auth.SetSuperUserContext(cli.Context())
				err error
				ntf = service.DefaultAuthNotification.With(ctx)
			)

			err = ntf.EmailConfirmation("en", args[0], "notification-testing-token")
			cli.HandleError(err)

			err = ntf.PasswordReset("en", args[0], "notification-testing-token")
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
