package provision

import (
	"context"
	"os"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

// Sets email-related settings (if not set) under "auth.internal..."
//
//
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
			_, err := store.LookupSettingByNameOwnedBy(ctx, s, setting.Name, 0)
			if errors.IsNotFound(err) {
				setting.UpdatedAt = *now()
				return store.CreateSetting(ctx, s, setting)
			}

			return err
		})
	})
}
