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

		rs := []string{ss[0].Name, ss[1].Name, ss[2].Name, ss[3].Name}
		req.Subset(rs, []string{"s1.opt.1", "s1.opt.2", "s2.opt.1", "s2.opt.2"})
	})
}
