package commands

import (
	"context"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/system/internal/auth/external"
	"github.com/cortezaproject/corteza-server/system/internal/repository"
	"github.com/cortezaproject/corteza-server/system/internal/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Will perform OpenID connect auto-configuration
func Auth(ctx context.Context, c *cli.Config) *cobra.Command {
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
			c.InitServices(ctx, c)

			var (
				name, providerUrl = args[0], args[1]
				es                = service.DefaultAuthSettings
			)

			if !skipValidationOnAutoDiscoveredProvider {
				// Do basic validation of external auth settings
				// will fail if secret or url are not set
				cli.HandleError(es.StaticValidateExternal())

				// Do full rediredct-URL check
				cli.HandleError(es.ValidateExternalRedirectURL())
			}

			p, err := parseExternalProviderUrl(providerUrl)
			cli.HandleError(err)

			eap, err := external.RegisterNewOpenIdClient(ctx, es, name, p.String())
			cli.HandleError(err)

			vv, err := eap.MakeValueSet("openid-connect." + name)
			cli.HandleError(err)

			if enableDiscoveredProvider {
				cli.HandleError(vv.Walk(func(value *settings.Value) error {
					if strings.HasSuffix(value.Name, ".enabled") {
						return value.SetValue(true)
					}

					return nil
				}))

				v := &settings.Value{Name: "auth.external.enabled"}
				cli.HandleError(v.SetValue(true))
				vv = append(vv, v)
			}

			cli.HandleError(service.DefaultIntSettings.BulkSet(vv))

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

			c.InitServices(ctx, c)

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
			c.InitServices(ctx, c)

			var (
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

func parseExternalProviderUrl(in string) (p *url.URL, err error) {
	if i := strings.Index(in, "://"); i == -1 {
		// Add schema if missing
		in = "https://" + in
	}

	if p, err = url.Parse(in); err != nil {
		// Try to parse it
		return
	} else if i := strings.Index(p.Path, external.WellKnown); i > -1 {
		// Cut off well-known-path
		p.Path = p.Path[:i]
	}

	return
}
