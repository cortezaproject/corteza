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
		_, err := (&Service{}).Bake(ctx, EncodeParams{Envoy: EnvoyConfig{MergeAlg: OnConflictPanic, SkipIf: "false"}}, nil, a, b, c)
		req.NoError(err)

		req.Equal(OnConflictPanic, a.Config.MergeAlg)
		req.Equal("false", a.Config.SkipIf)

		req.Equal(OnConflictReplace, b.Config.MergeAlg)
		req.Equal("false", b.Config.SkipIf)

		req.Equal(OnConflictReplace, c.Config.MergeAlg)
		req.Equal("true", c.Config.SkipIf)
	})

	t.Run("precompute expressions", func(t *testing.T) {
		a := &Node{
			Config: EnvoyConfig{
				MergeAlg: OnConflictDefault,
				SkipIf:   "",
			},
		}
		b := &Node{
			Config: EnvoyConfig{
				MergeAlg: OnConflictReplace,
				SkipIf:   "a == b",
			},
		}

		_, err := (&Service{}).Bake(ctx, EncodeParams{}, nil, a, b)
		req.NoError(err)

		req.Nil(a.Config.SkipIfEval)
		req.NotNil(b.Config.SkipIfEval)
	})
}
