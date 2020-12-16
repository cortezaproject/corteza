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
		bl := NewBuilder()

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
		bl := NewBuilder()

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
		bl := NewBuilder()

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

func TestGraph_WalkFlags(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("simple node link; a -> b -> c -> c", func(t *testing.T) {
		bl := NewBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id2": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id3": true}}},
		}
		c := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id3": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id3": true}}},
		}

		rr := []resource.Interface{
			a,
			b,
			c,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)

		s, err := g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(c, s.Res)
		req.True(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(b, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(c, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Nil(s)
	})

	t.Run("simple node link; a -> b -> c -> d -> b; a -> b -> e", func(t *testing.T) {
		bl := NewBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id2": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id3": true}}, &resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id5": true}}},
		}
		c := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id3": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id4": true}}},
		}
		d := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id4": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}
		e := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id5": true},
			refs:        nil,
		}

		rr := []resource.Interface{
			a,
			b,
			c,
			d,
			e,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)

		s, err := g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(e, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(b, s.Res)
		req.True(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(d, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(c, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(b, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Nil(s)
	})

	t.Run("simple node link; a -> b -> c -> d -> b; a -> b -> c -> e -> b", func(t *testing.T) {
		bl := NewBuilder()

		a := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id1": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		b := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id2": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id3": true}}},
		}
		c := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id3": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id4": true}}, &resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id5": true}}},
		}
		d := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id4": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}
		e := &testResource{
			resType:     "test:resource:1:",
			identifiers: resource.Identifiers{"id5": true},
			refs:        resource.RefSet{&resource.Ref{ResourceType: "test:resource:1:", Identifiers: resource.Identifiers{"id2": true}}},
		}

		rr := []resource.Interface{
			a,
			b,
			c,
			d,
			e,
		}

		g, err := bl.Build(ctx, rr...)
		req.NoError(err)

		s, err := g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(b, s.Res)
		req.True(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(a, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(d, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(e, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(c, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Equal(b, s.Res)
		req.False(s.Conflicting)

		s, err = g.NextInverted(ctx)
		req.NoError(err)
		req.Nil(s)
	})
}

func TestGraph_Utilities(t *testing.T) {
	req := require.New(t)

	t.Run("adding & removing nodes", func(t *testing.T) {
		g := newGraph()
		n1 := newNode(nil)
		n2 := newNode(nil)

		// Adding nodes
		g.addNode(n1, n2)
		req.Len(g.nn, 2)
		req.Equal(n1, g.nn[0])
		req.Equal(n2, g.nn[1])

		// Dependencies
		g.addChild(n1, n2)
		g.addParent(n2, n1)
		req.True(g.childNodes(n1).has(n2))
		req.True(g.parentNodes(n2).has(n1))

		// Removing nodes
		g.removeNode(n1, n2)
		req.False(g.childNodes(n1).has(n2))
		req.False(g.parentNodes(n2).has(n1))
	})

	t.Run("inverted graph", func(t *testing.T) {
		g := newGraph()
		n1 := newNode(nil)
		n2 := newNode(nil)

		// Adding nodes
		g.addNode(n1, n2)
		req.Len(g.nn, 2)
		req.Equal(n1, g.nn[0])
		req.Equal(n2, g.nn[1])

		// Dependencies
		g.addChild(n1, n2)
		g.addParent(n2, n1)
		g.invert()
		req.True(g.parentNodes(n1).has(n2))
		req.True(g.childNodes(n2).has(n1))
	})

	t.Run("inverted child dependencies", func(t *testing.T) {
		g := newGraph()
		n1 := newNode(nil)
		n2 := newNode(nil)
		n3 := newNode(nil)

		g.addNode(n1, n2, n3)

		g.addChild(n1, n2)
		g.invert()
		g.addChild(n1, n3)
		req.Equal(n2, n1.cc[0])
		req.Equal(n3, n1.pp[0])

		g.removeChild(n1, n3)
		req.Len(n1.pp, 0)
	})

	t.Run("inverted parent dependencies", func(t *testing.T) {
		g := newGraph()
		n1 := newNode(nil)
		n2 := newNode(nil)
		n3 := newNode(nil)

		g.addNode(n1, n2, n3)

		g.addParent(n1, n2)
		g.invert()
		g.addParent(n1, n3)
		req.Equal(n2, n1.pp[0])
		req.Equal(n3, n1.cc[0])

		g.removeParent(n1, n3)
		req.Len(n1.cc, 0)
	})

	// This one is just for codecoverage; the logic is nothing
	t.Run("remove all child/parent nodes", func(t *testing.T) {
		g := newGraph()
		n1 := newNode(nil)
		n2 := newNode(nil)

		g.addNode(n1, n2)
		g.addChild(n1, n2)
		g.addParent(n1, n2)

		g.removeChild(n1)
		g.removeParent(n1)
		req.Len(n1.cc, 0)
		req.Len(n1.pp, 0)

		// Inverted graph
		g.addChild(n1, n2)
		g.addParent(n1, n2)
		g.invert()

		g.removeChild(n1)
		g.removeParent(n1)
		req.Len(n1.cc, 0)
		req.Len(n1.pp, 0)
	})
}
