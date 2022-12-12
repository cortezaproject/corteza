package provision

import (
	"context"
	"fmt"
	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/mail"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/service"
	"go.uber.org/zap"
	"os"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
)

var (
	nextID = func() uint64 {
		return id.Next()
	}
)

// Sets email-related settings (if not set) under "auth.internal..."
func emailSettings(ctx context.Context, s store.Storer) error {
	var (
		val, has      = os.LookupEnv("SMTP_HOST")
		canSendEmails = has && len(val) > 0
	)

	// List of name-value pairs we need to iterate and set
	ss := types.SettingValueSet{
		types.MakeSettingValue(
			"auth.internal.signup.email-confirmation-required",
			canSendEmails,
		),

		types.MakeSettingValue(
			"auth.internal.password-reset.enabled",
			canSendEmails,
		),
	}

	return s.Tx(ctx, func(ctx context.Context, s store.Storer) error {
		return ss.Walk(func(setting *types.SettingValue) error {
			_, err := store.LookupSettingValueByNameOwnedBy(ctx, s, setting.Name, 0)
			if errors.IsNotFound(err) {
				setting.UpdatedAt = *now()
				return store.CreateSettingValue(ctx, s, setting)
			}

			return err
		})
	})
}

// addAuthSuperUsers assigns BYPASS roles to users from AUTH_PROVISION_SUPER_USER value
// When in Production, Corteza should stop and report an error.
func addAuthSuperUsers(ctx context.Context, log *zap.Logger, s store.Storer, authOpt options.AuthOpt) (err error) {
	var (
		envOpt = options.Environment()
	)

	if authOpt.ProvisionSuperUser == "" {
		return nil
	}

	if envOpt.IsProduction() {
		log.Warn(fmt.Sprint("when in production environment (ENVIRONMENT=production) you cannot provision " +
			"super users; set the environment to dev (ENVIRONMENT=dev) to provision super users"))
		return
	}

	users := strings.Split(authOpt.ProvisionSuperUser, ";")

	for _, usr := range users {
		u := prepareUser(usr)

		//check if the email address is valid
		if !mail.IsValidAddress(u.Email) {
			log.Warn(fmt.Sprintf("Email address %s is invalid", u.Email))
			continue
		}

		// skip existing email
		_, err = s.LookupUserByEmail(ctx, u.Email)
		if err != store.ErrNotFound {
			log.Warn(fmt.Sprintf("Email address already %s exists", u.Email))
			continue
		}

		// skip existing handle
		_, err = s.LookupUserByHandle(ctx, u.Handle)
		if err != store.ErrNotFound {
			log.Warn(fmt.Sprintf("Handle %s already exists", u.Handle))
			continue
		}

		err = store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
			if err = store.CreateUser(ctx, s, u); err != nil {
				return err
			}

			if err = service.SetPasswordCredentials(ctx, s, u.ID, u.Email); err != nil {
				return err
			}

			log.Warn(fmt.Sprintf("User {userID: %d, email: %s} created", u.ID, u.Email))

			//assign the user a bypass role
			for _, r := range internalAuth.BypassRoles() {
				m := &types.RoleMember{UserID: u.ID, RoleID: r.ID}
				if err = store.CreateRoleMember(ctx, s, m); err != nil {
					return err
				}
			}

			return
		})

		if err != nil {
			return err
		}
	}

	return
}

// prepareUser creates and fills a new user depending on the number of arguments provided
func prepareUser(user string) *types.User {
	u := &types.User{
		ID:        nextID(),
		CreatedAt: *now(),
	}

	usr := strings.Split(user, ",")
	u.Email = usr[0]
	u.EmailConfirmed = true
	u.Handle = createUserHandle(u)

	if len(usr) > 1 {
		u.Handle = usr[1]
	}

	if len(usr) > 2 {
		u.Name = usr[2]
	}

	return u
}

func createUserHandle(u *types.User) (hdl string) {
	hdl, _ = handle.Cast(
		func(lookup string) bool {
			return true
		},
		// use email w/o domain
		regexp.
			MustCompile("(@.*)$").
			ReplaceAllString(u.Email, ""),
	)

	return hdl
}
