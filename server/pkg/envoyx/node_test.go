package envoyx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeForRef(t *testing.T) {
	req := require.New(t)

	ref := Ref{
		ResourceType: "a",
		Identifiers:  MakeIdentifiers("a1", "a2"),
		Scope: Scope{
			ResourceType: "a",
			Identifiers:  MakeIdentifiers("a1", "a2"),
		},
	}

	t.Run("no nodes defined", func(t *testing.T) {
		req.Nil(NodeForRef(ref))
	})

	t.Run("not found", func(t *testing.T) {
		nn := NodeForRef(ref, NodeSet{{
			// The type misses
			ResourceType: "b",
			Identifiers:  MakeIdentifiers("a1", "a2"),
			Scope: Scope{
				ResourceType: "a",
				Identifiers:  MakeIdentifiers("a1", "a2"),
			},
		}}...)
		req.Nil(nn)
	})

	t.Run("not found wrong scope", func(t *testing.T) {
		nn := NodeForRef(ref, NodeSet{{
			// The type misses
			ResourceType: "a",
			Identifiers:  MakeIdentifiers("a1", "a2"),
			Scope: Scope{
				ResourceType: "b",
				Identifiers:  MakeIdentifiers("a1", "a2"),
			},
		}}...)
		req.Nil(nn)
	})

	t.Run("found", func(t *testing.T) {
		nn := NodeForRef(ref, NodeSet{{
			// The type misses
			ResourceType: "a",
			Identifiers:  MakeIdentifiers("a1", "a2"),
			Scope: Scope{
				ResourceType: "a",
				Identifiers:  MakeIdentifiers("a1", "a2"),
			},
		}}...)
		req.NotNil(nn)
	})
}

func TestIdents(t *testing.T) {
	req := require.New(t)

	ii := MakeIdentifiers("asdf", "asd123", "321das", "12367831412")

	ints, ss := ii.Idents()
	req.Len(ss, 3)
	req.Contains(ss, "asdf")
	req.Contains(ss, "asd123")
	req.Contains(ss, "321das")

	req.Len(ints, 1)
	req.Contains(ints, uint64(12367831412))
}

func TestIntersection(t *testing.T) {
	req := require.New(t)

	aa := MakeIdentifiers("a", "b")
	bb := MakeIdentifiers("a", "b")
	cc := MakeIdentifiers("b", "c")
	dd := MakeIdentifiers("c", "d")

	t.Run("complete overlap", func(t *testing.T) {
		req.True(aa.HasIntersection(bb))
	})

	t.Run("partial overlap", func(t *testing.T) {
		req.True(aa.HasIntersection(cc))
	})

	t.Run("completely off", func(t *testing.T) {
		req.False(aa.HasIntersection(dd))
	})
}

func TestFriendlyIdentifier(t *testing.T) {
	req := require.New(t)

	t.Run("empty", func(t *testing.T) {
		ii := MakeIdentifiers()
		req.Equal(ii.FriendlyIdentifier(), "")
	})

	t.Run("regular", func(t *testing.T) {
		ii := MakeIdentifiers("h1", "h2", "123123123")
		req.Equal(ii.FriendlyIdentifier(), "h1")
	})

	t.Run("fallback to ID", func(t *testing.T) {
		ii := MakeIdentifiers("123123123")
		req.Equal(ii.FriendlyIdentifier(), "123123123")
	})
}
