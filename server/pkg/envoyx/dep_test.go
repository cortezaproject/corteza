package envoyx

import (
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/stretchr/testify/require"
)

func TestScopeNodes(t *testing.T) {
	req := require.New(t)

	t.Run("none scoped", func(t *testing.T) {
		a := &Node{}
		b := &Node{}
		c := &Node{}

		ss := scopeNodes(a, b, c)
		req.Len(ss, 1)
		req.Len(ss[0], 3)
		req.Contains(ss[0], a)
	})

	t.Run("all same scope", func(t *testing.T) {
		a := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("a"),
			},
		}
		b := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("a"),
			},
		}
		c := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("a"),
			},
		}

		ss := scopeNodes(a, b, c)
		req.Len(ss, 1)
		req.Len(ss[0], 3)
		req.Contains(ss[0], a)
	})

	t.Run("mixed scope", func(t *testing.T) {
		none1 := &Node{}
		none2 := &Node{}

		a1 := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("a"),
			},
		}
		a2 := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("a"),
			},
		}

		b1 := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("b"),
			},
		}
		b2 := &Node{
			Scope: Scope{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  MakeIdentifiers("b"),
			},
		}

		ss := scopeNodes(none1, none2, a1, a2, b1, b2)
		req.Len(ss, 3)

		req.Len(ss[0], 2)
		req.Len(ss[1], 2)
		req.Len(ss[2], 2)

		req.True(ss[0][0].Scope.Identifiers.HasIntersection(ss[0][1].Scope.Identifiers))
		req.True(ss[1][0].Scope.Identifiers.HasIntersection(ss[1][1].Scope.Identifiers))
		req.True(ss[2][0].Scope.Identifiers.HasIntersection(ss[2][1].Scope.Identifiers))
	})
}

func TestDepRoots(t *testing.T) {
	req := require.New(t)

	t.Run("no nodes", func(t *testing.T) {
		gg := BuildDepGraph()
		req.Len(gg.Roots(), 0)
	})

	t.Run("single non-needy node", func(t *testing.T) {
		a1 := &Node{
			ResourceType: systemTypes.UserResourceType,
			Identifiers:  MakeIdentifiers("A.1"),
		}

		gg := BuildDepGraph(a1)

		req.Len(gg.Roots(), 1)
		req.Contains(gg.Roots(), a1)
	})

	t.Run("simple compose ns-mod dep setup", func(t *testing.T) {
		a1 := &Node{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("A.1"),
		}
		b1 := &Node{
			ResourceType: types.ModuleResourceType,
			Identifiers:  MakeIdentifiers("B.1"),
			References: map[string]Ref{
				"NamespaceID": a1.ToRef(),
			},
		}

		gg := BuildDepGraph(a1, b1)

		req.Len(gg.Roots(), 1)
		req.Contains(gg.Roots(), a1)
	})

	t.Run("simple compose mod-mod field dep setup", func(t *testing.T) {
		a1 := &Node{
			ResourceType: types.ModuleResourceType,
			Identifiers:  MakeIdentifiers("A.1"),
		}
		b1 := &Node{
			ResourceType: types.ModuleFieldResourceType,
			Identifiers:  MakeIdentifiers("B.1"),
			References: map[string]Ref{
				"NamespaceID": a1.ToRef(),
			},
		}

		gg := BuildDepGraph(a1, b1)

		req.Len(gg.Roots(), 1)
		req.Contains(gg.Roots(), a1)
	})
}

func TestDepTraversal(t *testing.T) {
	req := require.New(t)

	ns1 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns1"),
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns1"),
		},
	}
	ns2 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns2"),
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns2"),
		},
	}

	ns1mod1 := &Node{
		ResourceType: types.ModuleResourceType,
		Identifiers:  MakeIdentifiers("mod"),
		References: map[string]Ref{
			"NamespaceID": ns1.ToRef(),
		},
		Scope: ns1.Scope,
	}
	ns1mod1f1 := &Node{
		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  MakeIdentifiers("f1"),
		References: map[string]Ref{
			"NamespaceID": ns1.ToRef(),
			"ModuleID":    ns1mod1.ToRef(),
		},
		Scope: ns1.Scope,
	}

	ns2mod1 := &Node{
		ResourceType: types.ModuleResourceType,
		Identifiers:  MakeIdentifiers("mod"),
		References: map[string]Ref{
			"NamespaceID": ns2.ToRef(),
		},
		Scope: ns2.Scope,
	}
	ns2mod1f1 := &Node{
		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  MakeIdentifiers("f1"),
		References: map[string]Ref{
			"NamespaceID": ns2.ToRef(),
			"ModuleID":    ns2mod1.ToRef(),
		},
		Scope: ns2.Scope,
	}

	gg := BuildDepGraph(ns1, ns2, ns1mod1, ns1mod1f1, ns2mod1, ns2mod1f1)

	t.Run("children of namespace root", func(t *testing.T) {
		cc := gg.Children(ns1)

		req.Len(cc, 2)
		req.Contains(cc, ns1mod1)
		req.Contains(cc, ns1mod1f1)
	})

	t.Run("children module", func(t *testing.T) {
		cc := gg.Children(ns1mod1)

		req.Len(cc, 1)
		req.Contains(cc, ns1mod1f1)
	})

	t.Run("parent by ref missing", func(t *testing.T) {
		p := gg.ParentForRef(ns1mod1, Ref{ResourceType: types.NamespaceResourceType, Identifiers: MakeIdentifiers("asdf"), Scope: ns1.Scope})

		req.Nil(p)
	})

	t.Run("parent by ref wrong scope", func(t *testing.T) {
		p := gg.ParentForRef(ns1mod1, Ref{ResourceType: types.NamespaceResourceType, Identifiers: ns2.Identifiers, Scope: ns2.Scope})

		req.Nil(p)
	})

	t.Run("parent by ref", func(t *testing.T) {
		p := gg.ParentForRef(ns1mod1, Ref{ResourceType: types.NamespaceResourceType, Identifiers: ns1.Identifiers, Scope: ns1.Scope})

		req.NotNil(p)
		req.Equal(ns1, p)
	})
}

func TestDepXLinking(t *testing.T) {
	req := require.New(t)

	s1 := &Node{
		ResourceType: systemTypes.UserResourceType,
		Identifiers:  MakeIdentifiers("u1"),
	}
	s2 := &Node{
		ResourceType: systemTypes.UserResourceType,
		Identifiers:  MakeIdentifiers("u2"),
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns1"),
		},
	}

	ns1 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns1"),
		References: map[string]Ref{
			"S1": s1.ToRef(),
		},
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns1"),
		},
	}
	ns2 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns2"),
		References: map[string]Ref{
			"S2": s2.ToRef(),
		},
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns2"),
		},
	}
	_ = ns2

	t.Run("scoped node access unscoped ref", func(t *testing.T) {
		gg := BuildDepGraph(s1, s2, ns1)
		req.NotNil(gg.ParentForRef(ns1, s1.ToRef()))
	})

	t.Run("does not have missing refs", func(t *testing.T) {
		gg := BuildDepGraph(s1, s2, ns1)
		mm := gg.MissingRefs()
		req.Len(mm, 0)
	})

	t.Run("scoped node prevented ref of different scope", func(t *testing.T) {
		gg := BuildDepGraph(s1, s2, ns2)
		req.Nil(gg.ParentForRef(ns2, s2.ToRef()))
	})

	t.Run("has missing refs", func(t *testing.T) {
		gg := BuildDepGraph(s1, s2, ns2)
		mm := gg.MissingRefs()
		req.Len(mm, 1)
	})
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza/server/pkg/envoyx
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkDepGraphConstruction-12    	     826	   1495619 ns/op	 1633445 B/op	    6168 allocs/op
// PASS
func BenchmarkDepGraphConstruction(b *testing.B) {
	ns1 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns1"),
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns1"),
		},
	}
	ns2 := &Node{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  MakeIdentifiers("ns2"),
		Scope: Scope{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  MakeIdentifiers("ns2"),
		},
	}

	ns1mod1 := &Node{
		ResourceType: types.ModuleResourceType,
		Identifiers:  MakeIdentifiers("mod"),
		References: map[string]Ref{
			"NamespaceID": ns1.ToRef(),
		},
		Scope: ns1.Scope,
	}
	ns1mod1f1 := &Node{
		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  MakeIdentifiers("f1"),
		References: map[string]Ref{
			"NamespaceID": ns1.ToRef(),
			"ModuleID":    ns1mod1.ToRef(),
		},
		Scope: ns1.Scope,
	}

	ns2mod1 := &Node{
		ResourceType: types.ModuleResourceType,
		Identifiers:  MakeIdentifiers("mod"),
		References: map[string]Ref{
			"NamespaceID": ns2.ToRef(),
		},
		Scope: ns2.Scope,
	}
	ns2mod1f1 := &Node{
		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  MakeIdentifiers("f1"),
		References: map[string]Ref{
			"NamespaceID": ns2.ToRef(),
			"ModuleID":    ns2mod1.ToRef(),
		},
		Scope: ns2.Scope,
	}

	qwerty := make(NodeSet, 0, 1000-6)
	for i := 0; i < 1000-6; i++ {
		qwerty = append(qwerty, &Node{
			ResourceType: types.ModuleFieldResourceType,
			Identifiers:  MakeIdentifiers(fmt.Sprintf("gg_f_%d", i)),
			References: map[string]Ref{
				"NamespaceID": ns1.ToRef(),
				"ModuleID":    ns1mod1.ToRef(),
			},
			Scope: ns1.Scope,
		})
	}

	qwerty = append(qwerty, NodeSet{ns1, ns2, ns1mod1, ns1mod1f1, ns2mod1, ns2mod1f1}...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildDepGraph(qwerty...)
	}
}
