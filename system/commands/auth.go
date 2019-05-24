package commands

import (
	"context"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/system/internal/auth/external"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Will perform OpenID connect auto-configuration
func Auth(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "External authentication",
	}

	autoDiscoverCmd := &cobra.Command{
		Use:   "auto-discovery [name] [url]",
		Short: "Auto discovers new OIDC client",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var name, url = args[0], args[1]

			if eas, err := external.ExternalAuthSettings(service.DefaultIntSettings); err != nil {
				exit(err)
			} else if eap, err := external.RegisterNewOpenIdClient(ctx, eas, name, url); err != nil {
				exit(err)
			} else if vv, err := eap.MakeValueSet("openid-connect." + name); err != nil {
				exit(err)
			} else if err := service.DefaultIntSettings.BulkSet(vv); err != nil {
				exit(err)
			}
		},
	}

	jwtCmd := &cobra.Command{
		Use:   "jwt [email-or-id]",
		Short: "Generates new JWT for a user",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				db = factory.Database.MustGet("system")

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
				rr, err = roleRepo.FindByMemberID(user.ID)
			}

			if err != nil {
				exit(err)
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
				err error
				ntf = service.DefaultAuthNotification.With(ctx)
			)

			err = ntf.EmailConfirmation("en", args[0], "notification-testing-token")
			if err != nil {
				exit(err)
			}

			err = ntf.PasswordReset("en", args[0], "notification-testing-token")
			if err != nil {
				exit(err)
			}

		},
	}

	cmd.AddCommand(
		autoDiscoverCmd,
		testEmails,
		jwtCmd,
	)

	return cmd
}
