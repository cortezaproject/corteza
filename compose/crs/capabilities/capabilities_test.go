package capabilities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonCapabilities(t *testing.T) {
	cases := []struct {
		name string
		aa   Set
		bb   Set
		cc   Set
	}{{
		name: "regular",
		aa: Set{
			Create,
			Update,
		},
		bb: Set{
			Update,
			Delete,
		},
		cc: Set{
			Update,
		},
	}, {
		name: "no commons",
		aa: Set{
			Create,
		},
		bb: Set{
			Delete,
		},
		cc: Set{},
	}, {
		name: "aa empty",
		aa:   nil,
		bb: Set{
			Update,
			Delete,
		},
		cc: Set{},
	}, {
		name: "bb empty",
		aa: Set{
			Create,
			Update,
		},
		bb: nil,
		cc: Set{},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cc := common(c.aa, c.bb)
			require.Equal(t, c.cc, cc)
		})
	}
}

func TestCapabilityChecking(t *testing.T) {
	cases := []struct {
		name    string
		support Set
		require Set
		out     bool
	}{{
		name: "passing: complete match",
		support: Set{
			Create,
			Update,
		},
		require: Set{
			Create,
			Update,
		},
		out: true,
	}, {
		name: "passing: supports more then required",
		support: Set{
			Create,
			Update,
			Delete,
		},
		require: Set{
			Create,
			Update,
		},
		out: true,
	}, {
		name: "passing: no required",
		support: Set{
			Create,
			Update,
			Delete,
		},
		require: Set{},
		out:     true,
	}, {
		name:    "passing: no required nor supported",
		support: Set{},
		require: Set{},
		out:     true,
	}, {
		name: "failing: missing support",
		support: Set{
			Create,
			Update,
		},
		require: Set{
			Delete,
		},
		out: false,
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := c.support.IsSuperset(c.require...)
			require.Equal(t, c.out, out)

			out = c.require.IsSubset(c.support...)
			require.Equal(t, c.out, out)
		})
	}
}
