package envoyx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBake(t *testing.T) {
	ctx := context.Background()
	req := require.New(t)

	t.Run("set default envoy config", func(t *testing.T) {
		a := &Node{
			Config: EnvoyConfig{
				MergeAlg: OnConflictDefault,
				SkipIf:   "",
			},
		}
		b := &Node{
			Config: EnvoyConfig{
				MergeAlg: OnConflictReplace,
				SkipIf:   "",
			},
		}
		c := &Node{
			Config: EnvoyConfig{
				MergeAlg: OnConflictReplace,
				SkipIf:   "true",
			},
		}
		err := (&service{}).Bake(ctx, EncodeParams{Envoy: EnvoyConfig{MergeAlg: OnConflictPanic, SkipIf: "false"}}, a, b, c)
		req.NoError(err)

		req.Equal(OnConflictPanic, a.Config.MergeAlg)
		req.Equal("false", a.Config.SkipIf)

		req.Equal(OnConflictReplace, b.Config.MergeAlg)
		req.Equal("false", b.Config.SkipIf)

		req.Equal(OnConflictReplace, c.Config.MergeAlg)
		req.Equal("true", c.Config.SkipIf)
	})
}
