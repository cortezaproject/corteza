package dal

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/filter"
)

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza/server/pkg/dal
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkRowComparator-12    	  161826	      7015 ns/op	    5120 B/op	      40 allocs/op
// PASS
func BenchmarkRowComparator(b *testing.B) {
	cmp := makeRowComparator(
		&filter.SortExpr{Column: "a", Descending: false},
		&filter.SortExpr{Column: "b", Descending: true},
		&filter.SortExpr{Column: "c", Descending: true},
		&filter.SortExpr{Column: "d", Descending: false},
	)

	r1 := simpleRow{"a": 10, "b": "aa", "c": 33, "d": 500}
	r2 := simpleRow{"a": 50, "b": "aaaaa", "c": 31, "d": 10}
	r3 := simpleRow{"a": 31, "b": "a", "c": 11, "d": 1000}
	r4 := simpleRow{"a": 42, "b": "", "c": 0, "d": 300}
	r5 := simpleRow{"a": 22, "b": "aaaaaaaaaa", "c": -25, "d": 21}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmp(r1, r2)
		cmp(r1, r3)
		cmp(r1, r4)
		cmp(r1, r5)
		cmp(r2, r3)
		cmp(r2, r4)
		cmp(r2, r5)
		cmp(r3, r4)
		cmp(r3, r5)
		cmp(r4, r5)
	}
}
