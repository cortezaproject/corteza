package envoy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGraph_Walk(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("simple node link; a -> b", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: ResourceIdentifiers{"id1": true},
			refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: ResourceIdentifiers{"id2": true},
			refs:        nil,
		}

		exp := []Resource{b, a}
		rr := []Resource{
			a,
			b,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		i := 0
		g.Walk(ctx, func(ctx context.Context, s *ExecState) (NodeState, error) {
			req.Equal(exp[i], s.Res)
			i++
			return nil, nil
		})
	})

	t.Run("cyclic node link; a -> b -> a", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: ResourceIdentifiers{"id1": true},
			refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: ResourceIdentifiers{"id2": true},
			refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id1": true}}},
		}

		exp := []Resource{a, b, a}
		rr := []Resource{
			a,
			b,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		i := 0
		g.Walk(ctx, func(ctx context.Context, s *ExecState) (NodeState, error) {
			req.Equal(exp[i], s.Res)
			i++
			return nil, nil
		})
	})

	t.Run("self-cycle; a -> a", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: ResourceIdentifiers{"id1": true},
			refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id1": true}}},
		}

		exp := []Resource{a, a}
		rr := []Resource{
			a,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		g.invert()

		i := 0
		g.Walk(ctx, func(ctx context.Context, s *ExecState) (NodeState, error) {
			req.Equal(exp[i], s.Res)
			i++
			return nil, nil
		})
	})
}
