package expr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNxtRange(t *testing.T) {
	tcc := []struct {
		path string

		start  int
		expEnd int
		expErr bool
	}{{
		path: "a",

		start:  0,
		expEnd: 1,
	}, {
		path: "bcd",

		start:  0,
		expEnd: 3,
	}, {
		path: "bcd[e]",

		start:  0,
		expEnd: 3,
	}, {
		path: "bcd.e",

		start:  0,
		expEnd: 3,
	}, {
		path: "bcd.",

		start:  0,
		expErr: true,
	}, {
		path: "bcd[",

		start:  0,
		expErr: true,
	}, {
		path: "bcd.e",

		start:  4,
		expEnd: 5,
	}}

	for _, c := range tcc {
		t.Run(c.path[c.start:], func(t *testing.T) {
			_, o, _, err := nxtRange(c.path, c.start)
			if c.expErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, c.expEnd, o)
		})
	}
}

func TestPath(t *testing.T) {
	tcc := []struct {
		path string

		expBits  []string
		expRests []string
		expErr   bool
	}{{
		path: "a",

		expBits:  []string{"a"},
		expRests: []string{""},
	}, {
		path: "a.b",

		expBits:  []string{"a", "b"},
		expRests: []string{"b", ""},
	}, {
		path: "a.Content-Type",

		expBits:  []string{"a", "Content-Type"},
		expRests: []string{"Content-Type", ""},
	}, {
		path: "a[0]",

		expBits:  []string{"a", "0"},
		expRests: []string{"[0]", ""},
	}, {
		path: "a[b][c]",

		expBits:  []string{"a", "b", "c"},
		expRests: []string{"[b][c]", "[c]", ""},
	}, {
		path: "aa.b[c][d].e.f[g][0]",

		expBits:  []string{"aa", "b", "c", "d", "e", "f", "g", "0"},
		expRests: []string{"b[c][d].e.f[g][0]", "[c][d].e.f[g][0]", "[d].e.f[g][0]", "e.f[g][0]", "f[g][0]", "[g][0]", "[0]", ""},
	},
	}

	for _, c := range tcc {
		t.Run(c.path, func(t *testing.T) {
			pp := Path(c.path)
			var err error

			i := -1
			for {
				i++

				err = pp.Next()
				require.NoError(t, err)

				if !pp.More() {
					break
				}

				require.Equal(t, c.expBits[i], pp.Get())
				require.Equal(t, c.expRests[i], pp.Rest())
			}
			require.Equal(t, len(c.expBits), i)

			// for _, b := range c.expBits {

			// 	require.Equal(t, b, pp.Get())
			// }
			// pp, err = pp.Next()
			// require.NoError(t, err)

			// require.False(t, pp.More())
		})
	}
}

func BenchmarkPath(b *testing.B) {
	path := "aa.b[c][d].e.f[g][0]"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pp := Path(path)

		for {
			pp.Next()
			if !pp.More() {
				break
			}
		}
	}
}
