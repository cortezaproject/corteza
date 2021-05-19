package saml

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

type (
	ssFn           func(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error)
	certMockLoader struct {
	}

	mockStorer struct {
		ss func(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error)
	}
)

var (
	c = &types.SettingValue{Name: fmt.Sprintf("%s.cert", settingsPrefix)}
	k = &types.SettingValue{Name: fmt.Sprintf("%s.key", settingsPrefix)}
)

func Test_loadCertFromSettings(t *testing.T) {

	tcc := []struct {
		name   string
		err    error
		expect *Cert
		ss     ssFn
	}{
		{
			name:   "cert, key load success",
			expect: &Cert{Cert: []byte("CERT"), Key: []byte("KEY")},
			ss: func(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
				c.SetValue("CERT")
				k.SetValue("KEY")

				return types.SettingValueSet{c, k}, f, nil
			},
		},
		{
			name:   "no permissions to fetch settings",
			err:    errors.New("no permissions"),
			expect: &Cert{},
			ss: func(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
				return nil, types.SettingsFilter{}, errors.New("no permissions")
			},
		},
		{
			name:   "no cert, key in store",
			expect: &Cert{},
			ss: func(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
				return nil, types.SettingsFilter{}, nil
			},
		},
	}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				loader = CertStoreLoader{
					Storer: mockStorer{
						ss: tc.ss,
					},
				}
			)

			cc, err := loader.Load(context.Background())

			if tc.err != nil {
				req.Error(err)
			} else {
				req.NoError(err)
			}

			req.Equal(tc.expect, cc)
		})
	}

}

func (cm mockStorer) SearchSettings(ctx context.Context, f types.SettingsFilter) (types.SettingValueSet, types.SettingsFilter, error) {
	return cm.ss(ctx, f)
}
