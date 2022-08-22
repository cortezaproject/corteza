package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRelIndexBuffer(t *testing.T) {
	cc := newRelIndexBuffer("a", "b")

	cc.add((&row{}).WithValue("a", 0, 1).WithValue("b", 0, "a"))

	require.Len(t, cc.rows, 1)
	require.Equal(t, 1, cc.min["a"])
	require.Equal(t, 1, cc.max["a"])
	require.Equal(t, "a", cc.min["b"])
	require.Equal(t, "a", cc.max["b"])

	cc.add((&row{}).WithValue("a", 0, -1).WithValue("b", 0, "a"))
	require.Equal(t, -1, cc.min["a"])
	require.Equal(t, 1, cc.max["a"])
	require.Equal(t, "a", cc.min["b"])
	require.Equal(t, "a", cc.max["b"])

	cc.add((&row{}).WithValue("a", 0, 2).WithValue("b", 0, "aa"))
	require.Equal(t, -1, cc.min["a"])
	require.Equal(t, 2, cc.max["a"])
	require.Equal(t, "a", cc.min["b"])
	require.Equal(t, "aa", cc.max["b"])
}
