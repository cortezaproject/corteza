package app

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/auth/external"
	authSettings "github.com/cortezaproject/corteza-server/auth/settings"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/require"
)

type (
	authSettingsUpdaterMockedAuthService struct {
		settings *authSettings.Settings
	}
)

func (*authSettingsUpdaterMockedAuthService) MountHttpRoutes(chi.Router) {}
func (*authSettingsUpdaterMockedAuthService) Watch(context.Context)      {}
func (m *authSettingsUpdaterMockedAuthService) UpdateSettings(settings *authSettings.Settings) {
	m.settings = settings
}

func Test_updateAuthSettings(t *testing.T) {
	cc := []struct {
		name string
		in   types.AppSettings
		out  authSettings.Settings
	}{
		{
			"enabled internal auth",
			func() (s types.AppSettings) {
				s.Auth.Internal.Enabled = true
				return
			}(),
			func() (s authSettings.Settings) {
				s.LocalEnabled = true
				return
			}(),
		},
		{
			"add external providers",
			func() (s types.AppSettings) {
				s.Auth.External.Providers = []*types.ExternalAuthProvider{
					{ // skip
						Enabled: false,
						Handle:  "disabled",
						Key:     "key",
						Secret:  "sec",
					},
					{ // skip
						Enabled: true,
						Handle:  "malformed",
						Key:     "",
						Secret:  "",
					},
					{ // skip
						Enabled: true,
						Handle:  external.OIDC_PROVIDER_PREFIX + "without-issuer",
						Key:     "key",
						Secret:  "sec",
					},
					{ // add
						Enabled:   true,
						Handle:    external.OIDC_PROVIDER_PREFIX + "with-issuer",
						IssuerUrl: "issuer",
						Key:       "key",
						Secret:    "sec",
					},
					{ // add (w/o issuer)
						Enabled:   true,
						Handle:    "google",
						IssuerUrl: "issuer",
						Key:       "key",
						Secret:    "sec",
					},
					{ // add
						Enabled: true,
						Handle:  "github",
						Key:     "key",
						Secret:  "sec",
					},
					{ // skip
						Enabled: false,
						Handle:  "linkedin",
						Key:     "key",
						Secret:  "sec",
					},
				}
				return
			}(),
			func() (s authSettings.Settings) {
				s.Providers = []authSettings.Provider{
					{
						Handle:    external.OIDC_PROVIDER_PREFIX + "with-issuer",
						IssuerUrl: "issuer",
						Key:       "key",
						Secret:    "sec",
					},
					{
						Handle: "google",
						Key:    "key",
						Secret: "sec",
					},
					{
						Handle: "github",
						Key:    "key",
						Secret: "sec",
					},
				}
				return
			}(),
		},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			var (
				req = require.New(t)
				m   = &authSettingsUpdaterMockedAuthService{}
			)

			updateAuthSettings(m, &c.in)
			req.Equal(*m.settings, c.out)
		})
	}
}
