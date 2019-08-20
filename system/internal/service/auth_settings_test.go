package service

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx/types"

	intset "github.com/cortezaproject/corteza-server/internal/settings"
)

func Test_extractProviders(t *testing.T) {
	type args struct {
		redirectUrl string
		kv          intset.KV
	}
	tests := []struct {
		name          string
		args          args
		wantProviders map[string]AuthSettingsExternalAuthProvider
		wantErr       bool
	}{
		{
			name: "Empty KV",
			args: args{},
			wantProviders: map[string]AuthSettingsExternalAuthProvider{
				"github":   AuthSettingsExternalAuthProvider{},
				"linkedin": AuthSettingsExternalAuthProvider{},
				"google":   AuthSettingsExternalAuthProvider{},
				"facebook": AuthSettingsExternalAuthProvider{},
			},
		},

		{
			name: "Random config",
			args: args{
				kv: intset.KV{
					"auth.external.redirect-url":                         types.JSONText(`"http://%s"`),
					"auth.external.providers.openid-connect.foo.enabled": types.JSONText("true"),
					"auth.external.providers.openid-connect.foo.issuer":  types.JSONText(`"url"`),
					"auth.external.providers.openid-connect.foo.key":     types.JSONText(`"key"`),
					"auth.external.providers.openid-connect.foo.secret":  types.JSONText(`"secret"`),
					"auth.external.providers.openid-connect.bar.enabled": types.JSONText("true"),
					"auth.external.providers.openid-connect.bar.issuer":  types.JSONText(`"url"`),
					"auth.external.providers.openid-connect.bar.key":     types.JSONText(`"key"`),
					"auth.external.providers.openid-connect.bar.secret":  types.JSONText(`"secret"`),
					"auth.external.providers.openid-connect.baz.enabled": types.JSONText("false"),
					"auth.external.providers.github.enabled":             types.JSONText(`false`),
					"auth.external.providers.facebook.enabled":           types.JSONText(`true`),
					"auth.external.providers.facebook.secret":            types.JSONText(`"fb-secret"`),
				},
			},

			wantProviders: map[string]AuthSettingsExternalAuthProvider{
				"openid-connect.foo": AuthSettingsExternalAuthProvider{
					Enabled:     true,
					Key:         "key",
					Secret:      "secret",
					RedirectUrl: "http://openid-connect.foo",
					IssuerUrl:   "url",
				},
				"openid-connect.bar": AuthSettingsExternalAuthProvider{
					Enabled:     true,
					Key:         "key",
					Secret:      "secret",
					RedirectUrl: "http://openid-connect.bar",
					IssuerUrl:   "url",
				},
				"openid-connect.baz": AuthSettingsExternalAuthProvider{},
				"github":             AuthSettingsExternalAuthProvider{},
				"linkedin":           AuthSettingsExternalAuthProvider{},
				"google":             AuthSettingsExternalAuthProvider{},
				"facebook": AuthSettingsExternalAuthProvider{
					Enabled:     true,
					Secret:      "fb-secret",
					RedirectUrl: "http://facebook",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as, err := ParseAuthSettings(tt.args.kv)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(as.ExternalProviders, tt.wantProviders) {
				t.Errorf("extractProviders()\ngot:  %v\nwant: %v\n", as.ExternalProviders, tt.wantProviders)
			}
		})
	}
}
