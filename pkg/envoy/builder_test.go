package graph

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	testResource struct {
		resType     string
		identifiers ResourceIdentifiers
		refs        NodeRefSet
	}
)

func (t *testResource) Identifiers() ResourceIdentifiers {
	return t.identifiers
}

func (t *testResource) ResourceType() string {
	return t.resType
}

func (t *testResource) Refs() NodeRefSet {
	return t.refs
}

func TestGraphBuilder_Rel(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("single, simple node; a", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		rr := []Resource{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id1": true},
				refs:        nil,
			},
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Empty(g.childNodes(a))
		req.Empty(g.parentNodes(a))
	})

	t.Run("simple node link; a -> b", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		rr := []Resource{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id1": true},
				refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id2": true}}},
			},
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id2": true},
				refs:        nil,
			},
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		req.Len(g.nodes(), 2)

		a := g.nodes()[0]
		b := g.nodes()[1]
		req.Len(g.childNodes(a), 1)
		req.Equal(b, g.childNodes(a)[0])
		req.Empty(g.parentNodes(a))

		req.Len(b.pp, 1)
		req.Equal(a, b.pp[0])
	})

	t.Run("cyclic node link; a -> b -> a", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		rr := []Resource{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id1": true},
				refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id2": true}}},
			},
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id2": true},
				refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id1": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
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
		bl := NewGraphBuilder(nil, nil)

		rr := []Resource{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id1": true},
				refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id2": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Len(g.childNodes(a), 0)
		req.Len(g.parentNodes(a), 0)
	})

	t.Run("self-cycle; a -> a", func(t *testing.T) {
		bl := NewGraphBuilder(nil, nil)

		rr := []Resource{
			&testResource{
				resType:     "test:resource:1:",
				identifiers: ResourceIdentifiers{"id1": true},
				refs:        NodeRefSet{&NodeRef{ResourceType: "test:resource:1:", Identifiers: ResourceIdentifiers{"id1": true}}},
			},
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)
		req.Len(g.nodes(), 1)

		a := g.nodes()[0]
		req.Len(g.childNodes(a), 1)
		req.Equal(g.childNodes(a)[0], a)
		req.Len(g.parentNodes(a), 1)
		req.Equal(g.parentNodes(a)[0], a)
	})
}
