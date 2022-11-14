package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func pathHelper(pp ...string) (out [][]string) {
	for _, p := range pp {
		out = append(out, []string{p})
	}
	return
}

func TestSimplePath(t *testing.T) {
	n := NewIndex()

	n.Add(1, pathHelper("a")...)

	out := n.Collect(pathHelper("a")...)
	require.Len(t, out, 1)
	require.Equal(t, 1, out[0])
}

func TestSimpleLongPath(t *testing.T) {
	n := NewIndex()

	n.Add(1, pathHelper("a", "b", "c")...)

	out := n.Collect(pathHelper("a", "b", "c")...)
	require.Len(t, out, 1)
	require.Equal(t, 1, out[0])
}

func TestWildSimplePath(t *testing.T) {
	n := NewIndex()

	n.Add(1, pathHelper("a", "*")...)

	out := n.Collect(pathHelper("a", "b")...)
	require.Len(t, out, 1)
	require.Equal(t, 1, out[0])
}

func TestComplex(t *testing.T) {
	n := NewIndex()

	n.Add(1, pathHelper("resource", "lvl1.id1", "lvl2.id2", "res.id")...)
	n.Add(2, pathHelper("resource", "lvl1.id1", "lvl2.id2", "*")...)
	n.Add(3, pathHelper("resource", "lvl1.id1", "lvl2.id2", "*")...)
	n.Add(4, pathHelper("resource", "*", "*", "*")...)

	var out []interface{}

	t.Run("full path", func(t *testing.T) {
		out = n.Collect(pathHelper("resource", "lvl1.id1", "lvl2.id2", "res.id")...)
		require.Len(t, out, 4)
		require.Equal(t, []interface{}{1, 2, 3, 4}, out)
	})

	t.Run("full path; id not matching", func(t *testing.T) {
		out = n.Collect(pathHelper("resource", "lvl1.id1", "lvl2.id2", "invalid")...)
		require.Len(t, out, 3)
		require.Equal(t, []interface{}{2, 3, 4}, out)
	})

	t.Run("partial matching path 1", func(t *testing.T) {
		out = n.Collect(pathHelper("resource", "invald", "invalid", "invalid")...)
		require.Len(t, out, 1)
		require.Equal(t, []interface{}{4}, out)
	})

	t.Run("not matching path", func(t *testing.T) {
		out = n.Collect(pathHelper("invalid", "invald", "invalid", "invalid")...)
		require.Len(t, out, 0)
	})
}
