package rbac

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// goos: darwin
// goarch: arm64
// pkg: github.com/crusttech/human/server/pkg/rbac
// cpu: Apple M3 Pro
// BenchmarkAccess_100_10_1000-12                   1552434               755.3 ns/op           320 B/op         11 allocs/op
// BenchmarkAccess_1000_100_10000-12                1314985              1053 ns/op             320 B/op         11 allocs/op
// BenchmarkAccess_10000_1000_10000-12               246471              4589 ns/op             780 B/op         13 allocs/op
// BenchmarkAccess_100000_10000_100000-12              2286            509592 ns/op            1119 B/op         16 allocs/op
func BenchmarkAccess_100_10_1000(b *testing.B) {
	benchmarkAccess(b, 100, 10, 1000)
}

func BenchmarkAccess_1000_100_10000(b *testing.B) {
	benchmarkAccess(b, 1000, 100, 10000)
}

func BenchmarkAccess_10000_1000_10000(b *testing.B) {
	benchmarkAccess(b, 10000, 1000, 10000)
}

func BenchmarkAccess_100000_10000_100000(b *testing.B) {
	benchmarkAccess(b, 100000, 10000, 10000)
}

func benchmarkAccess(b *testing.B, l1n, l2n, l3n int) {
	r := &groupNode{
		id:     id.MustNumID(1),
		handle: "root",
	}
	l1 := makeNodeSet(1, l1n)
	l2 := makeNodeSet(2, l2n)
	l3 := makeNodeSet(3, l3n)

	connectNodeSets([]*groupNode{r}, l1)
	connectNodeSets(l1, l2)
	connectNodeSets(l2, l3)
	connectNodeSets(l1, l3)

	x := append([]*groupNode{r}, l1...)
	x = append(x, l2...)
	x = append(x, l3...)

	root, _, err := buildOrgTree(x...)
	require.NoError(b, err)

	_ = root

	mm := []GroupMembers{}
	for _, x := range x {
		mm = append(mm, GroupMembers{
			group:   x,
			members: []id.ID{x.id},
		})
	}

	svc := &orgTree{logger: zap.NewNop()}
	err = svc.Rebuild(mm...)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		svc.IsAbove(pickRandomID(l2), pickRandomID(l1))
		svc.IsAbove(pickRandomID(l3), pickRandomID(l1))
		svc.IsAbove(pickRandomID(l3), pickRandomID(l2))
		svc.IsAbove(pickRandomID(l3), root.id)
	}
}

func pickRandomID(slice []*groupNode) id.ID {
	return slice[rand.Intn(len(slice))].id
}

func makeNodeSet(lvl int, n int) (out []*groupNode) {
	for i := 0; i < n; i++ {
		out = append(out, &groupNode{
			id:     id.MustNumID(uint64(lvl*1000 + i)),
			handle: fmt.Sprintf("%d_%d", lvl, i),
		})
	}

	return
}

func connectNodeSets(parents, children []*groupNode) {
	for i, c := range children {
		c.paths = append(c.paths, mkPp(parents[i%len(parents)].id)...)
	}
}

func BenchmarkBuildOrgTree(b *testing.B) {
	bits := []GroupMembers{
		{
			group: &groupNode{
				id:     id.MustNumID(1),
				handle: "1",
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(2),
				handle: "2",
				paths:  mkPp(id.MustNumID(1)),
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(3),
				handle: "3",
				paths:  mkPp(id.MustNumID(1)),
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(4),
				handle: "4",
				paths:  mkPp(id.MustNumID(1)),
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(5),
				handle: "5",
				paths:  mkPp(id.MustNumID(1)),
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(6),
				handle: "6",
				paths:  mkPp(id.MustNumID(1)),
			},
		},

		{
			group: &groupNode{
				id:     id.MustNumID(7),
				handle: "7",
				paths:  mkPp(id.MustNumID(2)),
			},
		},
		{
			group: &groupNode{
				id:     id.MustNumID(8),
				handle: "8",
				paths:  mkPp(id.MustNumID(3)),
			},
		},
	}

	svc := &orgTree{logger: zap.NewNop()}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := svc.Rebuild(bits...)
		if err != nil {
			panic(err)
		}
	}
}
