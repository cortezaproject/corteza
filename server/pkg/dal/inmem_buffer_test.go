package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func TestInmemBuffer_rw(t *testing.T) {
	gBuff := InMemoryBuffer()
	ctx := context.Background()

	tcc := []struct {
		name string
		prep func() *inmemBuffer
		in   []simpleRow
		test func(t *testing.T, in []simpleRow, buff *inmemBuffer)
	}{{
		name: "fresh insert",
		prep: func() *inmemBuffer {
			return gBuff
		},
		in: []simpleRow{
			{
				"r1_k1": "r1_k1",
				"r1_k2": "r1_k2",
				"r1_k3": "r1_k3",
			}},
		test: func(t *testing.T, in []simpleRow, buff *inmemBuffer) {
			r := make(simpleRow)

			require.True(t, buff.Next(ctx))
			require.NoError(t, buff.Scan(r))
			require.Equal(t, in[0], r)

			require.False(t, buff.Next(ctx))
		},
	}, {
		name: "insert into existing",
		prep: func() *inmemBuffer {
			return gBuff
		},
		in: []simpleRow{
			{
				"r2_k1": "r2_k1",
				"r2_k2": "r2_k2",
				"r2_k3": "r2_k3",
			}},
		test: func(t *testing.T, in []simpleRow, buff *inmemBuffer) {
			r := make(simpleRow)

			require.True(t, buff.Next(ctx))
			require.NoError(t, buff.Scan(r))
			require.Equal(t, in[0], r)

			require.False(t, buff.Next(ctx))
		},
	}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			buff := c.prep()
			for _, in := range c.in {
				require.NoError(t, buff.Add(ctx, in))
			}
			c.test(t, c.in, buff)
		})
	}
}

func TestInmemBuffer_single_rw(t *testing.T) {
	mk := func() *inmemBuffer {
		buff := InMemoryBuffer()
		buff.Single()
		return buff
	}
	ctx := context.Background()

	tcc := []struct {
		name string
		prep func() *inmemBuffer
		in   []simpleRow
		test func(t *testing.T, in []simpleRow, buff *inmemBuffer)
	}{{
		name: "one",
		prep: func() *inmemBuffer {
			return mk()
		},
		in: []simpleRow{
			{
				"r1_k1": "r1_k1",
				"r1_k2": "r1_k2",
				"r1_k3": "r1_k3",
			}},
		test: func(t *testing.T, in []simpleRow, buff *inmemBuffer) {
			r := make(simpleRow)

			require.True(t, buff.Next(ctx))
			require.NoError(t, buff.Scan(r))
			require.Equal(t, in[0], r)

			require.False(t, buff.Next(ctx))
		},
	}, {
		name: "two",
		prep: func() *inmemBuffer {
			return mk()
		},
		in: []simpleRow{
			{
				"r1_k1": "r1_k1",
				"r1_k2": "r1_k2",
				"r1_k3": "r1_k3",
			}, {
				"r2_k1": "r2_k1",
				"r2_k2": "r2_k2",
				"r2_k3": "r2_k3",
			}},
		test: func(t *testing.T, in []simpleRow, buff *inmemBuffer) {
			r := make(simpleRow)

			require.True(t, buff.Next(ctx))
			require.NoError(t, buff.Scan(r))
			require.Equal(t, in[1], r)

			require.False(t, buff.Next(ctx))
		},
	}}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			buff := c.prep()
			for _, in := range c.in {
				require.NoError(t, buff.Add(ctx, in))
			}
			c.test(t, c.in, buff)
		})
	}
}

func TestInmemBuffer_single_writeMidRead(t *testing.T) {
	ctx := context.Background()

	buff := InMemoryBuffer()
	buff.Single()

	a := simpleRow{
		"r1_k1": "r1_k1",
		"r1_k2": "r1_k2",
		"r1_k3": "r1_k3",
	}

	b := simpleRow{
		"r2_k1": "r2_k1",
		"r2_k2": "r2_k2",
		"r2_k3": "r2_k3",
	}

	c := simpleRow{
		"r3_k1": "r3_k1",
		"r3_k2": "r3_k2",
		"r3_k3": "r3_k3",
	}

	// Write
	require.NoError(t, buff.Add(ctx, a))

	// Read
	tmp := simpleRow{}
	require.True(t, buff.Next(ctx))
	require.NoError(t, buff.Err())
	require.NoError(t, buff.Scan(tmp))
	require.Equal(t, a, tmp)

	// Next write
	require.NoError(t, buff.Add(ctx, b))
	// Next write
	require.NoError(t, buff.Add(ctx, c))

	// Read
	tmp = simpleRow{}
	require.True(t, buff.Next(ctx))
	require.NoError(t, buff.Err())
	require.NoError(t, buff.Scan(tmp))
	require.Equal(t, c, tmp)

	// Read empty
	require.False(t, buff.Next(ctx))
	require.NoError(t, buff.Err())
}

func TestInmemBuffer_inOrder(t *testing.T) {
	mk := func() *inmemBuffer {
		buff := InMemoryBuffer()
		return buff
	}
	ctx := context.Background()

	tcc := []struct {
		name string
		prep func() *inmemBuffer
		in   []simpleRow
		out  []simpleRow
		sort filter.SortExprSet
	}{{
		name: "in order",
		prep: func() *inmemBuffer {
			return mk()
		},
		in: []simpleRow{
			{
				"order": 3,
			}, {
				"order": 2,
			}, {
				"order": 4,
			}, {
				"order": 1,
			}},
		out: []simpleRow{
			{
				"order": 1,
			}, {
				"order": 2,
			}, {
				"order": 3,
			}, {
				"order": 4,
			}},
		sort: filter.SortExprSet{{
			Column:     "order",
			Descending: false,
		}},
	},
	}

	for _, c := range tcc {
		t.Run(c.name, func(t *testing.T) {
			buff := c.prep()
			buff.InOrder(c.sort...)
			for _, in := range c.in {
				require.NoError(t, buff.Add(ctx, in))
			}

			out := make([]simpleRow, 0, 4)
			for buff.Next(ctx) {
				require.NoError(t, buff.Err())

				r := simpleRow{}
				require.NoError(t, buff.Scan(r))
				out = append(out, r)
			}

			require.Equal(t, c.out, out)
		})
	}
}

func TestInmemBuffer_seek(t *testing.T) {
	ctx := context.Background()

	buff := InMemoryBuffer()

	a := simpleRow{
		"k": "k1",
	}
	b := simpleRow{
		"k": "k2",
	}
	c := simpleRow{
		"k": "k3",
	}

	tmp := simpleRow{}

	// Write
	require.NoError(t, buff.Add(ctx, a))
	require.NoError(t, buff.Add(ctx, b))
	require.NoError(t, buff.Add(ctx, c))

	// Seek to end
	require.NoError(t, buff.Seek(ctx, 4))
	require.False(t, buff.Next(ctx))

	// Seek to start
	require.NoError(t, buff.Seek(ctx, 0))
	require.True(t, buff.Next(ctx))
	require.NoError(t, buff.Scan(tmp))
	require.Equal(t, a, tmp)

	// Seek to middle
	require.NoError(t, buff.Seek(ctx, 2))
	require.True(t, buff.Next(ctx))
	require.NoError(t, buff.Scan(tmp))
	require.Equal(t, c, tmp)
}
