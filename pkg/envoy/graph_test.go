package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/stretchr/testify/require"
)

func TestGraph_Walk(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("simple node link; a -> b", func(t *testing.T) {
		bl := NewGraphBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id2": true},
			refs:        nil,
		}

		rr := []resource.Interface{
			a,
			b,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		s, err := g.Next(ctx)
		req.Equal(b, s.Res)
		req.False(s.Conflicting)

		s, err = g.Next(ctx)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)
	})

	t.Run("cyclic node link; a -> b -> a", func(t *testing.T) {
		bl := NewGraphBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id2": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id1": true}}},
		}

		rr := []resource.Interface{
			a,
			b,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		s, err := g.Next(ctx)
		req.Equal(a, s.Res)
		req.True(s.Conflicting)

		s, err = g.Next(ctx)
		req.Equal(b, s.Res)
		req.False(s.Conflicting)

		s, err = g.Next(ctx)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)
	})

	t.Run("self-cycle; a -> a", func(t *testing.T) {
		bl := NewGraphBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id1": true}}},
		}

		rr := []resource.Interface{
			a,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		s, err := g.Next(ctx)
		req.Equal(a, s.Res)
		req.True(s.Conflicting)

		s, err = g.Next(ctx)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)
	})
}
