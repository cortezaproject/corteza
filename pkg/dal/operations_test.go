package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonOperations(t *testing.T) {
	cases := []struct {
		name string
		aa   OperationSet
		bb   OperationSet
		cc   OperationSet
	}{{
		name: "regular",
		aa: OperationSet{
			Create,
			Update,
		},
		bb: OperationSet{
			Update,
			Delete,
		},
		cc: OperationSet{
			Update,
		},
	}, {
		name: "no commons",
		aa: OperationSet{
			Create,
		},
		bb: OperationSet{
			Delete,
		},
		cc: OperationSet{},
	}, {
		name: "aa empty",
		aa:   nil,
		bb: OperationSet{
			Update,
			Delete,
		},
		cc: OperationSet{},
	}, {
		name: "bb empty",
		aa: OperationSet{
			Create,
			Update,
		},
		bb: nil,
		cc: OperationSet{},
	}}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cc := common(c.aa, c.bb)
			require.Equal(t, c.cc, cc)
		})
	}
}

func TestOperationChecking(t *testing.T) {
	cases := []struct {
		name    string
		support OperationSet
		require OperationSet
		out     bool
	}{{
		name: "passing: complete match",
		support: OperationSet{
			Create,
			Update,
		},
		require: OperationSet{
			Create,
			Update,
		},
		out: true,
	}, {
		name: "passing: supports more then required",
		support: OperationSet{
			Create,
			Update,
			Delete,
		},
		require: OperationSet{
			Create,
			Update,
		},
		out: true,
	}, {
		name: "passing: no required",
		support: OperationSet{
			Create,
			Update,
			Delete,
		},
		require: OperationSet{},
		out:     true,
	}, {
		name:    "passing: no required nor supported",
		support: OperationSet{},
		require: OperationSet{},
		out:     true,
	}, {
		name: "failing: missing support",
		support: OperationSet{
			Create,
			Update,
		},
		require: OperationSet{
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
