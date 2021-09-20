package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	"github.com/stretchr/testify/require"
)

type (
	testResource struct {
		resType     string
		identifiers resource.Identifiers
		refs        resource.RefSet
		ph          bool
	}
)

func (t *testResource) Identifiers() resource.Identifiers {
	return t.identifiers
}

func (t *testResource) ResourceType() string {
	return t.resType
}

func (t *testResource) Refs() resource.RefSet {
	return t.refs
}

func (t *testResource) MarkPlaceholder() {
	t.ph = true
}

func (t *testResource) Placeholder() bool {
	return t.ph
}

func TestGraphBuilder_Rel(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("single, simple node; a", func(t *testing.T) {
		bl := NewBuilder()

		rr := []resource.Interface{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id1": true},
				refs:        nil,
			},
		}

		g, err := bl.Build(ctx, rr...)
		g.invert()
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Empty(g.childNodes(a))
		req.Empty(g.parentNodes(a))
	})

	t.Run("simple node link; a -> b", func(t *testing.T) {
		bl := NewBuilder()

		rr := []resource.Interface{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id1": true},
				refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
			},
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id2": true},
				refs:        nil,
			},
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		req.Len(g.nodes(), 2)

		a := g.resIndex[rr[0]]
		b := g.resIndex[rr[1]]
		req.Len(g.childNodes(a), 1)
		req.Equal(b, g.childNodes(a)[0])
		req.Empty(g.parentNodes(a))

		req.Len(g.parentNodes(b), 1)
		req.Equal(a, g.parentNodes(b)[0])
	})

	t.Run("cyclic node link; a -> b -> a", func(t *testing.T) {
		bl := NewBuilder()

		rr := []resource.Interface{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id1": true},
				refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
			},
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id2": true},
				refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id1": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
		g.invert()
		req.NoError(err)
		req.Len(g.nodes(), 2)

		a := g.nodes()[0]
		b := g.nodes()[1]
		req.Len(g.childNodes(a), 1)
		req.Equal(b, g.childNodes(a)[0])
		req.Len(g.parentNodes(a), 1)
		req.Equal(b, g.parentNodes(a)[0])

		req.Len(b.cc, 1)
		req.Equal(a, b.cc[0])
		req.Len(b.pp, 1)
		req.Equal(a, b.pp[0])
	})

	t.Run("node with missing dep; a -> nill", func(t *testing.T) {
		bl := NewBuilder()

		rr := []resource.Interface{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id1": true},
				refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
		g.invert()
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Len(g.childNodes(a), 0)
		req.Len(g.parentNodes(a), 0)
	})

	t.Run("self-cycle; a -> a", func(t *testing.T) {
		bl := NewBuilder()

		rr := []resource.Interface{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: resource.Identifiers{"id1": true},
				refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id1": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
		g.invert()
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Len(g.childNodes(a), 1)
		req.Equal(g.childNodes(a)[0], a)
		req.Len(g.parentNodes(a), 1)
		req.Equal(g.parentNodes(a)[0], a)
	})
}
