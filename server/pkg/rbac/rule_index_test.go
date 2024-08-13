package rbac

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndexBuild(t *testing.T) {
	tcc := []struct {
		name   string
		in     []*Rule
		remove []*Rule
		add    []*Rule
		out    []int

		role uint64
		op   string
		res  string
	}{
		{
			name: "empty",
			in:   nil,
			out:  nil,

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "match",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "read",
				Access:    Allow,
			}},
			out: []int{0},

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
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

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "one match one role missmatch",
			in: []*Rule{{
				RoleID:    2,
				Resource:  "a:b/c/d",
				Operation: "read",
				Access:    Allow,
			}, {
				RoleID:    1,
				Resource:  "a:b/*/*",
				Operation: "read",
				Access:    Inherit,
			}},
			out: []int{1},

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "role missmatch",
			in: []*Rule{{
				RoleID:    2,
				Resource:  "a:b/c/d",
				Operation: "read",
				Access:    Allow,
			}, {
				RoleID:    3,
				Resource:  "a:b/*/*",
				Operation: "read",
				Access:    Inherit,
			}},
			out: nil,

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "path missmatch",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/e",
				Operation: "read",
				Access:    Allow,
			}},
			out: nil,

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "operation missmatch",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: nil,

			role: 1,
			op:   "read",
			res:  "a:b/c/d",
		}, {
			name: "removing the only element",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			remove: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: nil,

			role: 1,
			op:   "write",
			res:  "a:b/c/d",
		},
		{
			name: "removing twice added thing",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}, {
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			remove: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},

			out: nil,

			role: 1,
			op:   "write",
			res:  "a:b/c/d",
		},
		{
			name: "two elements with no common root",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}, {
				RoleID:    2,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			remove: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: nil,

			role: 1,
			op:   "write",
			res:  "a:b/c/d",
		},
		{
			name: "two elements with common root (get removed)",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}, {
				RoleID:    1,
				Resource:  "a:b/c/e",
				Operation: "write",
				Access:    Allow,
			}},
			remove: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: nil,

			role: 1,
			op:   "write",
			res:  "a:b/c/d",
		},
		{
			name: "two elements with common root (get not removed)",
			in: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}, {
				RoleID:    1,
				Resource:  "a:b/c/e",
				Operation: "write",
				Access:    Allow,
			}},
			remove: []*Rule{{
				RoleID:    1,
				Resource:  "a:b/c/d",
				Operation: "write",
				Access:    Allow,
			}},
			out: []int{1},

			role: 1,
			op:   "write",
			res:  "a:b/c/e",
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

			role: 1,
			op:   "write",
			res:  "a:b/c/x",
		}}

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			ix := buildRuleIndex(tc.in)
			ix.remove(tc.remove...)
			ix.add(tc.add...)

			out := RuleSet(ix.get(tc.role, tc.op, tc.res))
			sort.Sort(out)

			want := RuleSet(graby(append(tc.in, tc.add...), tc.out))
			sort.Sort(want)

			require.Len(t, out, len(want))
			for i, o := range out {
				require.Equal(t, want[i], o)
			}
		})
	}

}

func graby(rr []*Rule, want []int) (out []*Rule) {
	out = make([]*Rule, 0, len(want))

	for _, w := range want {
		out = append(out, rr[w])
	}

	return
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza/server/pkg/rbac
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkIndexBuild_100-12                 10000            102361 ns/op           88064 B/op       1271 allocs/op
// BenchmarkIndexBuild_1000-12                 1149           1024872 ns/op          755375 B/op      11183 allocs/op
// BenchmarkIndexBuild_10000-12                 128           8986248 ns/op         4406477 B/op      82453 allocs/op
// BenchmarkIndexBuild_100000-12                 14          81871407 ns/op        20627785 B/op     543568 allocs/op
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
