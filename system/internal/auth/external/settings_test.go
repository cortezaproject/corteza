package external

import (
	"reflect"
	"testing"

	intset "github.com/cortezaproject/corteza-server/internal/settings"
	"github.com/jmoiron/sqlx/types"
)

func Test_extractProviders(t *testing.T) {
	type args struct {
		redirectUrl string
		kv          intset.KV
	}
	tests := []struct {
		name          string
		args          args
		wantProviders map[string]externalAuthProvider
		wantErr       bool
	}{
		{
			name: "Empty KV",
			args: args{},
			wantProviders: map[string]externalAuthProvider{
				"github":   externalAuthProvider{},
				"linkedin": externalAuthProvider{},
				"gplus":    externalAuthProvider{},
				"facebook": externalAuthProvider{},
			},
		},

		{
			name: "Random config",
			args: args{
				redirectUrl: "http://%s",
				kv: intset.KV{
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

			wantProviders: map[string]externalAuthProvider{
				"openid-connect.foo": externalAuthProvider{
					enabled:     true,
					key:         "key",
					secret:      "secret",
					redirectUrl: "http://openid-connect.foo",
					issuerUrl:   "url",
				},
				"openid-connect.bar": externalAuthProvider{
					enabled:     true,
					key:         "key",
					secret:      "secret",
					redirectUrl: "http://openid-connect.bar",
					issuerUrl:   "url",
				},
				"openid-connect.baz": externalAuthProvider{
					enabled: false,
				},
				"github": externalAuthProvider{
					enabled: false,
				},
				"linkedin": externalAuthProvider{
					enabled: false,
				},
				"gplus": externalAuthProvider{
					enabled: false,
				},
				"facebook": externalAuthProvider{
					enabled:     true,
					secret:      "fb-secret",
					redirectUrl: "http://facebook",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProviders, err := extractProviders(tt.args.redirectUrl, tt.args.kv)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractProviders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotProviders, tt.wantProviders) {
				t.Errorf("extractProviders() = %v, want %v", gotProviders, tt.wantProviders)
			}
		})
	}
}

func TestExternalAuthProvider(t *testing.T) {
	type args struct {
		kv intset.KV
	}
	tests := []struct {
		name    string
		args    args
		wantEap externalAuthProvider
		wantErr bool
	}{
		{
			args: args{
				kv: intset.KV{
					"foo.enabled": types.JSONText("true"),
					"foo.issuer":  types.JSONText(`"example.tld"`),
					"foo.key":     types.JSONText(`"key"`),
					"foo.secret":  types.JSONText(`"secret"`),
				},
			},

			wantEap: externalAuthProvider{
				enabled:   true,
				key:       "key",
				secret:    "secret",
				issuerUrl: "example.tld",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEap, err := ExternalAuthProvider(tt.args.kv)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExternalAuthProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotEap, tt.wantEap) {
				t.Errorf("ExternalAuthProvider() = %v, want %v", gotEap, tt.wantEap)
			}
		})
	}
}
