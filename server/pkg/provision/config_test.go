package provision

import (
	"errors"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_collectPartialDirectories(t *testing.T) {
	var (
		dirs = []uConfig{{dir: "000_base"}, {dir: "003_auth"}, {dir: "300_automation"}}

		all = func(uConfig) bool { return true }

		present = map[string]envoyx.NodeSet{
			"000_base":       {&envoyx.Node{ResourceType: "base"}},
			"300_automation": {&envoyx.Node{ResourceType: "automation"}},
		}

		decode = func(dir string) (envoyx.NodeSet, error) { return present[dir], nil }
	)

	t.Run("a missing directory does not stop the ones after it", func(t *testing.T) {
		nn, err := collectPartialDirectories(zap.NewNop(), dirs, all, decode)
		require.NoError(t, err)
		require.Len(t, nn, 2)
		require.Equal(t, "automation", nn[1].ResourceType)
	})

	t.Run("a directory that needs no import is skipped", func(t *testing.T) {
		nn, err := collectPartialDirectories(zap.NewNop(), dirs, func(d uConfig) bool { return d.dir != "000_base" }, decode)
		require.NoError(t, err)
		require.Len(t, nn, 1)
		require.Equal(t, "automation", nn[0].ResourceType)
	})

	t.Run("a decode error names its cause", func(t *testing.T) {
		cause := errors.New("bad yaml")
		_, err := collectPartialDirectories(zap.NewNop(), dirs, all, func(string) (envoyx.NodeSet, error) { return nil, cause })
		require.ErrorIs(t, err, cause)
	})
}
