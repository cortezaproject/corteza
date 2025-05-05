package slice

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCycleSlice(t *testing.T) {
	req := require.New(t)

	cc := NewCircular[int](5)

	cc.Add(1)
	req.Equal([]int{1}, cc.Slice())

	cc.Add(2)
	cc.Add(3)
	req.Equal([]int{1, 2, 3}, cc.Slice())

	cc.Add(4)
	cc.Add(5)
	req.Equal([]int{1, 2, 3, 4, 5}, cc.Slice())

	cc.Add(6)
	req.Equal([]int{2, 3, 4, 5, 6}, cc.Slice())

	cc.Add(7)
	cc.Add(8)
	cc.Add(9)
	cc.Add(10)
	req.Equal([]int{6, 7, 8, 9, 10}, cc.Slice())

	cc.Add(11)
	cc.Add(12)
	req.Equal([]int{8, 9, 10, 11, 12}, cc.Slice())
}
