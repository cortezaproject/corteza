package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
	"github.com/stretchr/testify/require"
)

func TestSettings(t *testing.T) {
	var (
		ctx    = context.Background()
		s, err = initStore(ctx)
	)
	if err != nil {
		t.Fatalf("failed to init sqlite in-memory db: %v", err)
	}

	prepare := func(ctx context.Context, s store.Storer, t *testing.T, suite string) (*require.Assertions, error) {
		req := require.New(t)

		nn, err := dd(ctx, suite)
		req.NoError(err)

		return req, encode(ctx, s, nn)
	}

	// Prepare
	s, err = initStore(ctx)
	err = ce(
		err,

		s.TruncateSettings(ctx),
	)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := prepare(ctx, s, t, "settings")
	req.NoError(err)

	t.Run("settings", func(t *testing.T) {
		ss, _, err := store.SearchSettings(ctx, s, types.SettingsFilter{})
		req.NoError(err)
		req.NotNil(ss)
		req.Len(ss, 4)

		req.Equal("s1.opt.1", ss[0].Name)
		req.Equal("s1.opt.2", ss[1].Name)
		req.Equal("s2.opt.1", ss[2].Name)
		req.Equal("s2.opt.2", ss[3].Name)
	})
}
