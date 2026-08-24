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
			"providers.nylas.key":                  sqlTypes.JSONText(`"nylas"`),

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

	require.Equal(t,
		aux.Providers.FindByHandle("nylas"),
		&ExternalAuthProvider{Enabled: false, Key: "nylas", Handle: "nylas", Label: "Nylas"})
}

func Test_appSettingsLogoDefaults(t *testing.T) {
	var (
		none        = AppSettings{}
		customLight = AppSettings{}
		customBoth  = AppSettings{}
	)

	customLight.UI.MainLogo = "attachment:1"
	customBoth.UI.MainLogo = "attachment:1"
	customBoth.UI.MainLogoDark = "attachment:2"

	d := none.WithDefaults()
	require.Equal(t, "/assets/logo.svg", d.UI.MainLogo)
	require.Equal(t, "/assets/logo-dark.svg", d.UI.MainLogoDark)
	require.Equal(t, "/assets/icon.svg", d.UI.IconLogo)
	require.Empty(t, d.UI.IconLogoDark)

	// a custom light logo serves the dark theme too, never the built-in dark one
	d = customLight.WithDefaults()
	require.Equal(t, "attachment:1", d.UI.MainLogo)
	require.Empty(t, d.UI.MainLogoDark)

	d = customBoth.WithDefaults()
	require.Equal(t, "attachment:2", d.UI.MainLogoDark)
}

func Test_settingsConnectionOAuthAppsDecode(t *testing.T) {
	type Dst struct {
		Apps ConnectionOAuthAppSet
	}

	var (
		aux = Dst{}
		kv  = SettingsKV{
			"apps.google.client-id":     sqlTypes.JSONText(`"g-cid"`),
			"apps.google.client-secret": sqlTypes.JSONText(`"g-secret"`),
			"apps.github.client-id":     sqlTypes.JSONText(`"gh-cid"`),
		}
	)

	require.NoError(t, DecodeKV(kv, &aux))
	require.Len(t, aux.Apps, 2)

	require.Equal(t,
		&ConnectionOAuthApp{Handle: "google", ClientID: "g-cid", ClientSecret: "g-secret"},
		aux.Apps.FindByHandle("google"))

	require.Equal(t,
		&ConnectionOAuthApp{Handle: "github", ClientID: "gh-cid"},
		aux.Apps.FindByHandle("github"))

	require.Nil(t, aux.Apps.FindByHandle("missing"))
}

func Test_settingsConnectionOAuthAppAccessor(t *testing.T) {
	var (
		as AppSettings
		kv = SettingsKV{
			"connection.oauth.redirect-url":              sqlTypes.JSONText(`"https://app/callback"`),
			"connection.oauth.apps.google.client-id":     sqlTypes.JSONText(`"cid"`),
			"connection.oauth.apps.google.client-secret": sqlTypes.JSONText(`"sec"`),
			"connection.oauth.apps.google.auth-url":      sqlTypes.JSONText(`"https://accounts.google.com/o/oauth2/v2/auth"`),
			"connection.oauth.apps.google.token-url":     sqlTypes.JSONText(`"https://oauth2.googleapis.com/token"`),
			"connection.oauth.apps.google.pkce":          sqlTypes.JSONText(`true`),
		}
	)

	require.NoError(t, DecodeKV(kv, &as))

	app, redirect, ok := as.ConnectionOAuthApp("google")
	require.True(t, ok)
	require.Equal(t, "https://app/callback", redirect)
	require.Equal(t, "cid", app.ClientID)
	require.Equal(t, "sec", app.ClientSecret)
	require.Equal(t, "https://oauth2.googleapis.com/token", app.TokenURL)
	require.True(t, app.PKCE)

	_, _, ok = as.ConnectionOAuthApp("missing")
	require.False(t, ok)
}
