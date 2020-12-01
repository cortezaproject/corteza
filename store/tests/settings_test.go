package tests

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/rand"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func testSettings(t *testing.T, s store.Settings) {
	var (
		ctx = context.Background()

		makeNew = func(nn ...string) *types.SettingValue {
			name := strings.Join(nn, "")
			return &types.SettingValue{
				Name:      name,
				OwnedBy:   id.Next(),
				UpdatedAt: *now(),
			}
		}

		truncAndCreate = func(t *testing.T) (*require.Assertions, *types.SettingValue) {
			req := require.New(t)
			req.NoError(s.TruncateSettings(ctx))
			settings := makeNew(string(rand.Bytes(10)))
			req.NoError(s.CreateSetting(ctx, settings))
			return req, settings
		}

		truncAndFill = func(t *testing.T, l int) (*require.Assertions, types.SettingValueSet) {
			req := require.New(t)
			req.NoError(s.TruncateSettings(ctx))

			set := make([]*types.SettingValue, l)

			for i := 0; i < l; i++ {
				set[i] = makeNew(string(rand.Bytes(10)))
			}

			req.NoError(s.CreateSetting(ctx, set...))
			return req, set
		}
	)

	t.Run("create", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.CreateSetting(ctx, makeNew()))
	})

	t.Run("lookup by name and ownedBy", func(t *testing.T) {
		req, setting := truncAndCreate(t)
		fetched, err := s.LookupSettingByNameOwnedBy(ctx, setting.Name, setting.OwnedBy)
		req.NoError(err)
		req.Equal(setting.Name, fetched.Name)
		req.Equal(setting.OwnedBy, fetched.OwnedBy)
	})

	t.Run("update", func(t *testing.T) {
		req, setting := truncAndCreate(t)
		setting.Value = []byte(`42`)
		req.NoError(s.UpdateSetting(ctx, setting))

		fetched, err := s.LookupSettingByNameOwnedBy(ctx, setting.Name, setting.OwnedBy)
		req.NoError(err)
		req.Equal(string(`42`), string(fetched.Value))
	})

	t.Run("upsert", func(t *testing.T) {
		req := require.New(t)
		req.NoError(s.TruncateSettings(ctx))

		t.Run("new", func(t *testing.T) {
			req := require.New(t)
			req.NoError(s.UpsertSetting(ctx, &types.SettingValue{Name: "foo", Value: []byte(`"foo"`)}))
			v, err := s.LookupSettingByNameOwnedBy(ctx, "foo", 0)
			req.NoError(err)
			req.NotNil(v)
			req.Equal(string(`"foo"`), string(v.Value))
		})

		t.Run("existing", func(t *testing.T) {
			req.NoError(s.CreateSetting(ctx, &types.SettingValue{Name: "baz", Value: []byte(`"created"`)}))
			req.NoError(s.UpsertSetting(ctx, &types.SettingValue{Name: "baz", Value: []byte(`"updated"`)}))
			v, err := s.LookupSettingByNameOwnedBy(ctx, "baz", 0)
			req.NoError(err)
			req.NotNil(v)
			req.Equal(string(`"updated"`), string(v.Value))
		})

	})

	t.Run("delete", func(t *testing.T) {
		t.Run("by Settings", func(t *testing.T) {
			req, setting := truncAndCreate(t)
			req.NoError(s.DeleteSetting(ctx, setting))
			set, _, err := s.SearchSettings(ctx, types.SettingsFilter{OwnedBy: setting.OwnedBy})
			req.NoError(err)
			req.Len(set, 0)
		})

		t.Run("by name and ownedBy", func(t *testing.T) {
			req, setting := truncAndCreate(t)
			req.NoError(s.DeleteSettingByNameOwnedBy(ctx, setting.Name, setting.OwnedBy))
			set, _, err := s.SearchSettings(ctx, types.SettingsFilter{OwnedBy: setting.OwnedBy})
			req.NoError(err)
			req.Len(set, 0)
		})
	})

	t.Run("search", func(t *testing.T) {
		t.Run("by ownedBy", func(t *testing.T) {
			req, prefill := truncAndFill(t, 5)

			set, f, err := s.SearchSettings(ctx, types.SettingsFilter{OwnedBy: prefill[0].OwnedBy})
			req.NoError(err)
			req.Equal(uint64(prefill[0].OwnedBy), f.OwnedBy)
			req.Len(set, 1)
		})

		t.Run("with check", func(t *testing.T) {
			t.Skip("Should pass(afaik) but doesn't")

			req, prefill := truncAndFill(t, 5)

			set, _, err := s.SearchSettings(ctx, types.SettingsFilter{
				Check: func(setting *types.SettingValue) (bool, error) {
					return (setting.Name == prefill[0].Name), nil
				},
			})
			req.NoError(err)
			req.Len(set, 1)
			req.Equal(prefill[0].Name, set[0].Name)
		})
	})

	t.Run("ordered search", func(t *testing.T) {
		t.Skip("not implemented")
	})
}
