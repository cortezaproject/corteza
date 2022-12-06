package provision

import (
	"context"
	"os"

	"github.com/cortezaproject/corteza/server/pkg/errors"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/pkg/rand"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/cortezaproject/corteza/server/system/types"
	"go.uber.org/zap"
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

// defaultAuthClient checks if default client exists (handle = AUTH_DEFAULT_CLIENT) and adds it
func defaultAuthClient(ctx context.Context, log *zap.Logger, s store.AuthClients, authOpt options.AuthOpt) error {
	if authOpt.DefaultClient == "" {
		// Default client not set
		return nil
	}

	c := &types.AuthClient{
		ID:     id.Next(),
		Handle: authOpt.DefaultClient,
		Meta: &types.AuthClientMeta{
			Name: "Corteza Web Applications",
		},
		ValidGrant: "authorization_code",
		RedirectURI: func() string {
			// Disabling protection by redirection URL for now, it caused too much confusion on simple setups
			//baseURL, _ := url.Parse(authOpt.BaseURL)
			//return fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Hostname())
			return ""
		}(),

		Secret:    string(rand.Bytes(64)),
		Scope:     "profile api",
		Enabled:   true,
		Trusted:   true,
		Security:  &types.AuthClientSecurity{},
		Labels:    nil,
		CreatedAt: *now(),
	}

	_, err := store.LookupAuthClientByHandle(ctx, s, c.Handle)
	if err == nil {
		return nil
	}

	if !errors.IsNotFound(err) {
		return err
	}

	if err = store.CreateAuthClient(ctx, s, c); err != nil {
		return err
	}

	log.Info(
		"Added OAuth2 client",
		zap.String("name", c.Meta.Name),
		zap.String("redirectURI", c.RedirectURI),
		zap.Uint64("clientId", c.ID),
	)

	return nil
}
