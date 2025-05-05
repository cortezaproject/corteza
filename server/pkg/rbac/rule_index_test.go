package rbac

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexBuild(t *testing.T) {
	tcc := []struct {
		name string
		in   []*Rule
		add  []*Rule
		out  []int

		op  string
		res string
	}{
		{
			name: "empty",
			in:   nil,
			out:  nil,

			op:  "read",
			res: "a:b/c/d",
		}, {
			name: "match",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "read",
				Access:    Allow,
			}},
			out: []int{0},

			op:  "read",
			res: "a:b/c/d",
		}, {
			name: "multiple matches",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "read",
				Access:    Allow,
			}, {
				RoleID:    1,
				Resource:  "a:b/*/*",
				Operation: "read",
				Access:    Inherit,
			}},
			out: []int{0, 1},

			op:  "read",
			res: "a:b/c/d",
		}, {
			name: "path missmatch",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/e",
				Operation: "read",
				Access:    Allow,
			}},
			out: nil,

			op:  "read",
			res: "a:b/c/d",
		}, {
			name: "operation missmatch",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: nil,

			op:  "read",
			res: "a:b/c/d",
		},
		{
			name: "add new element",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			add: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/x",
				Operation: "write",
				Access:    Allow,
			}},

			out: []int{1},

			op:  "write",
			res: "a:b/c/x",
		}}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			ix := buildRuleIndex(tc.in)
			ix.add(tc.add...)

			out := RuleSet(ix.get(tc.op, tc.res))
			sort.Sort(out)

			want := RuleSet(grabIndexMatches(append(tc.in, tc.add...), tc.out))
			sort.Sort(want)

			require.Len(t, out, len(want))
			for i, o := range out {
				require.Equal(t, want[i], o)
			}
		})
	}
}

func TestIndexHas(t *testing.T) {
	ix := buildRuleIndex([]*Rule{{
		RoleID:    1,
		Resource:  "a:b/c/x",
		Operation: "write",
		Access:    Allow,
	}})

	require.True(t, ix.has(&Rule{
		RoleID:    1,
		Resource:  "a:b/c/x",
		Operation: "write",
		Access:    Allow,
	}))

	require.False(t, ix.has(&Rule{
		RoleID:    1,
		Resource:  "a:b/c/*",
		Operation: "write",
		Access:    Allow,
	}))

	require.False(t, ix.has(&Rule{
		RoleID:    1,
		Resource:  "a:b/c/zz",
		Operation: "write",
		Access:    Allow,
	}))

}

func grabIndexMatches(rr []*Rule, want []int) (out []*Rule) {
	out = make([]*Rule, 0, len(want))

	for _, w := range want {
		out = append(out, rr[w])
	}

	return
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// cpu: Apple M3 Pro
// BenchmarkIndexBuild_100-12                140071              7721 ns/op           10936 B/op         19 allocs/op
// BenchmarkIndexBuild_1000-12                14124             83847 ns/op          121899 B/op         45 allocs/op
// BenchmarkIndexBuild_10000-12                1161            944820 ns/op         1461061 B/op        429 allocs/op
// BenchmarkIndexBuild_100000-12                 98          12371088 ns/op        13424712 B/op       3409 allocs/op
func benchmarkIndexBuild(b *testing.B, rules []*Rule) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buildRuleIndex(rules)
	}
}

func BenchmarkIndexBuild_100(b *testing.B) {
	benchmarkIndexBuild(b, makeRuleSet(100, 10))
}

func BenchmarkIndexBuild_1000(b *testing.B) {
	benchmarkIndexBuild(b, makeRuleSet(1000, 10))
}

func BenchmarkIndexBuild_10000(b *testing.B) {
	benchmarkIndexBuild(b, makeRuleSet(10000, 10))
}

func BenchmarkIndexBuild_100000(b *testing.B) {
	benchmarkIndexBuild(b, makeRuleSet(100000, 10))
}

func makeRuleSet(count int, roleCount int) (out RuleSet) {
	for i := 0; i < count; i++ {
		out = append(out, &Rule{
			RoleID:    uint64(1000 + int(rand.Intn(roleCount))),
			Resource:  randomResource(),
			Operation: randomOperation(),
			Access:    randomAccess(),
		})
	}

	return
}
