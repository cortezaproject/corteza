package types

import (
	"testing"

	sqlTypes "github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

// 	Hello! This file is auto-generated.

func Test_settingsExtAuthProvidersValidConfiguration(t *testing.T) {

	var (
		empty        = ExternalAuthProvider{}
		google       = ExternalAuthProvider{Enabled: true, Handle: "google", Key: "some-guid", Secret: "s3cret"}
		noIssuerOIDC = ExternalAuthProvider{Enabled: true, Handle: "openid-connect.bar", Key: "some-guid", Secret: "s3cret"}
		goodOIDC     = ExternalAuthProvider{Enabled: true, Handle: "openid-connect.bar", Key: "some-guid", Secret: "s3cret", IssuerUrl: "https://example.org"}
	)

	require.False(t, noIssuerOIDC.ValidConfiguration())
	require.True(t, goodOIDC.ValidConfiguration())
	require.False(t, empty.ValidConfiguration())
	require.True(t, google.ValidConfiguration())
}

func Test_settingsExtAuthProvidersDecode(t *testing.T) {
	type (
		Dst struct {
			Providers ExternalAuthProviderSet
		}
	)

	var (
		aux = Dst{
			Providers: ExternalAuthProviderSet{
				{Handle: "github"},
				{Handle: "nylas"},
				{Handle: "facebook"},
				{Enabled: true, Key: "g00gl3", Handle: "google"},
				{Handle: "linkedin"},
				{Enabled: true, Key: "K3Y", Handle: "openid-connect.remove"},
			},
		}
		kv = SettingsKV{
			"providers.foo.enabled":                sqlTypes.JSONText(`true`),
			"providers.openid-connect.bar.enabled": sqlTypes.JSONText(`true`),
			"providers.openid-connect.bar.key":     sqlTypes.JSONText(`"K3Y"`),
			"providers.google.enabled":             sqlTypes.JSONText(`true`),
			"providers.google.key":                 sqlTypes.JSONText(`"g00gl3"`),

			// Values with null should not be added!
			"providers.openid-connect.null.handle": sqlTypes.JSONText(`null`),

			// Values with null should not be added!
			"providers.openid-connect.remove.handle": sqlTypes.JSONText(`null`),
			"providers.openid-connect.remove.key":    sqlTypes.JSONText(`null`),
		}
	)

	require.NoError(t, DecodeKV(kv, &aux))
	require.Len(t, aux.Providers, 6)

	require.Nil(t,
		aux.Providers.FindByHandle("foo"))

	require.Equal(t,
		aux.Providers.FindByHandle("openid-connect.bar"),
		&ExternalAuthProvider{Enabled: true, Key: "K3Y", Handle: "openid-connect.bar", Label: "Bar"})

	require.Equal(t,
		aux.Providers.FindByHandle("google"),
		&ExternalAuthProvider{Enabled: true, Key: "g00gl3", Handle: "google", Label: "Google"})

	require.Equal(t,
		aux.Providers.FindByHandle("linkedin"),
		&ExternalAuthProvider{Handle: "linkedin", Label: "LinkedIn"})

	require.Equal(t,
		aux.Providers.FindByHandle("github"),
		&ExternalAuthProvider{Handle: "github", Label: "GitHub"})

	require.Equal(t,
		aux.Providers.FindByHandle("facebook"),
		&ExternalAuthProvider{Handle: "facebook", Label: "Facebook"})

}
