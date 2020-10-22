package envoy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	TestNode struct {
		rr NodeRelationships
		ii NodeIdentifiers
	}
)

func (n *TestNode) Identifiers() NodeIdentifiers {
	return n.ii
}

func (n *TestNode) Matches(resource string, identifiers ...string) bool {
	if resource != n.Resource() {
		return false
	}

	return n.Identifiers().HasAny(identifiers...)
}

func (n *TestNode) Resource() string {
	return "envoy:test:"
}

func (n *TestNode) Relations() NodeRelationships {
	return n.rr
}

func TestEnvoyGraph_Rel(t *testing.T) {
	req := require.New(t)

	t.Run("simple node, no rels", func(t *testing.T) {
		g := NewGraph()

		rr := NodeRelationships{}
		ii := NodeIdentifiers{"p1"}
		n := &TestNode{rr: rr, ii: ii}
		g.Add(n)

		cc := g.Children(n)
		req.Empty(cc)

		pp := g.Parents(n)
		req.Empty(pp)
	})

	t.Run("node with child and parent nodes", func(t *testing.T) {
		g := NewGraph()

		// The child node
		rr1 := NodeRelationships{}
		ii1 := NodeIdentifiers{"c1"}
		n1 := &TestNode{rr: rr1, ii: ii1}
		g.Add(n1)

		// The middle node
		rr2 := NodeRelationships{"envoy:test:": NodeIdentifiers{"c1"}}
		ii2 := NodeIdentifiers{"n"}
		n := &TestNode{rr: rr2, ii: ii2}
		g.Add(n)

		// The parent node
		rr3 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n"}}
		ii3 := NodeIdentifiers{"p1"}
		n3 := &TestNode{rr: rr3, ii: ii3}
		g.Add(n3)

		cc := g.Children(n)
		req.Len(cc, 1)
		req.Equal(n1, cc[0])

		pp := g.Parents(n)
		req.Len(pp, 1)
		req.Equal(n3, pp[0])
	})

	t.Run("(inverted) node with child and parent nodes", func(t *testing.T) {
		g := NewGraph()
		g.Invert()

		// The child node
		rr1 := NodeRelationships{}
		ii1 := NodeIdentifiers{"c1"}
		n1 := &TestNode{rr: rr1, ii: ii1}
		g.Add(n1)

		// The middle node
		rr2 := NodeRelationships{"envoy:test:": NodeIdentifiers{"c1"}}
		ii2 := NodeIdentifiers{"n"}
		n := &TestNode{rr: rr2, ii: ii2}
		g.Add(n)

		// The parent node
		rr3 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n"}}
		ii3 := NodeIdentifiers{"p1"}
		n3 := &TestNode{rr: rr3, ii: ii3}
		g.Add(n3)

		cc := g.Children(n)
		req.Len(cc, 1)
		req.Equal(n3, cc[0])

		pp := g.Parents(n)
		req.Len(pp, 1)
		req.Equal(n1, pp[0])
	})
}

func TestEnvoyGraph_DepResolution(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("simple acyclic linear graph; (n1) => (n2) => (n3)", func(t *testing.T) {
		g := NewGraph()

		rr1 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n2"}}
		ii1 := NodeIdentifiers{"n1"}
		n1 := &TestNode{rr: rr1, ii: ii1}
		g.Add(n1)

		rr2 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n3"}}
		ii2 := NodeIdentifiers{"n2"}
		n2 := &TestNode{rr: rr2, ii: ii2}
		g.Add(n2)

		rr3 := NodeRelationships{}
		ii3 := NodeIdentifiers{"n3"}
		n3 := &TestNode{rr: rr3, ii: ii3}
		g.Add(n3)

		// 1. n1 since it has no parent nodes
		n, pp, cc, err := g.Next(ctx)
		req.Equal(n1, n)
		req.NoError(err)
		req.NotEmpty(cc)
		req.Empty(pp)

		// 2. n2 since its parents are resolved
		n, pp, cc, err = g.Next(ctx)
		req.Equal(n2, n)
		req.NoError(err)
		req.NotEmpty(cc)
		req.NotEmpty(pp)

		// 2. n3 since its parents are resolved
		n, pp, cc, err = g.Next(ctx)
		req.Equal(n3, n)
		req.NoError(err)
		req.Empty(cc)
		req.NotEmpty(pp)
	})
}

func TestEnvoyGraph_GarbageCollector(t *testing.T) {
	req := require.New(t)
	ctx := context.Background()

	t.Run("simple acyclic linear graph; (n1) => (n2) => (n3)", func(t *testing.T) {
		g := NewGraph()

		rr1 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n2"}}
		ii1 := NodeIdentifiers{"n1"}
		n1 := &TestNode{rr: rr1, ii: ii1}
		g.Add(n1)

		rr2 := NodeRelationships{"envoy:test:": NodeIdentifiers{"n3"}}
		ii2 := NodeIdentifiers{"n2"}
		n2 := &TestNode{rr: rr2, ii: ii2}
		g.Add(n2)

		rr3 := NodeRelationships{}
		ii3 := NodeIdentifiers{"n3"}
		n3 := &TestNode{rr: rr3, ii: ii3}
		g.Add(n3)

		// n1 marked as processed; no garbage
		g.Next(ctx)
		req.Len(g.nodes, 3)
		req.Len(g.processed, 1)
		req.Equal(g.processed[0], n1)

		// n2 marked as processed; garbage in nodes, processed
		g.Next(ctx)
		req.Len(g.nodes, 2)
		req.Len(g.processed, 1)
		req.Equal(g.processed[0], n2)

		// all nodes processed; resetting graph
		g.Next(ctx)
		req.Len(g.nodes, 0)
		req.Len(g.processed, 0)
	})
}
