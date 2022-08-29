package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeRows(t *testing.T) {
	tcc := []struct {
		name    string
		a       *Row
		b       *Row
		mapping []AttributeMapping
		out     *Row
	}{{
		name: "full merge; no mapping",
		a:    (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, "hi").WithValue("attr2", 1, "hello"),
		b:    (&Row{}).WithValue("attr3", 0, true).WithValue("attr4", 0, "ee").WithValue("attr4", 1, 25),
		out:  (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, "hi").WithValue("attr2", 1, "hello").WithValue("attr3", 0, true).WithValue("attr4", 0, "ee").WithValue("attr4", 1, 25),
	}, {
		name: "full merge; no mapping; collision",
		a:    (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, "hi").WithValue("attr2", 1, "hello"),
		b:    (&Row{}).WithValue("attr2", 0, true).WithValue("attr3", 0, "ee").WithValue("attr3", 1, 25),
		out:  (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, true).WithValue("attr3", 0, "ee").WithValue("attr3", 1, 25),
	},

		{
			name: "mapped merge",
			a:    (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, "hi").WithValue("attr2", 1, "hello"),
			b:    (&Row{}).WithValue("attr3", 0, true).WithValue("attr4", 0, "ee").WithValue("attr4", 1, 25),
			out:  (&Row{}).WithValue("a", 0, 10).WithValue("b", 0, "hi").WithValue("b", 1, "hello").WithValue("c", 0, true).WithValue("d", 0, "ee").WithValue("d", 1, 25),
			mapping: saToMapping([]simpleAttribute{{
				ident:  "a",
				source: "attr1",
			}, {
				ident:  "b",
				source: "attr2",
			}, {
				ident:  "c",
				source: "attr3",
			}, {
				ident:  "d",
				source: "attr4",
			}}...),
		}, {
			name: "mapped merge with conflicts",
			a:    (&Row{}).WithValue("attr1", 0, 10).WithValue("attr2", 0, "hi").WithValue("attr2", 1, "hello"),
			b:    (&Row{}).WithValue("attr3", 0, true).WithValue("attr4", 0, "ee").WithValue("attr4", 1, 25),
			out:  (&Row{}).WithValue("a", 0, 10).WithValue("b", 0, true).WithValue("c", 0, "ee").WithValue("c", 1, 25),
			mapping: saToMapping([]simpleAttribute{{
				ident:  "a",
				source: "attr1",
			}, {
				ident:  "b",
				source: "attr2",
			}, {
				ident:  "b",
				source: "attr3",
			}, {
				ident:  "c",
				source: "attr4",
			}}...),
		}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			out := &Row{}
			err := mergeRows(c.mapping, out, c.a, c.b)
			require.NoError(t, err)
			require.Equal(t, c.out, out)
		})
	}
}
