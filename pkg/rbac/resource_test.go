package rbac

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResourceType(t *testing.T) {
	var (
		tcc = []struct {
			in  string
			exp string
		}{
			{"a:b/c/d", "a:b"},
			{"a:b/c", "a:b"},
			{"a:b/", "a:b"},
			{"a:b", "a:b"},
			{"a/", "a"},
			{"a", "a"},
			{"", ""},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.exp, ResourceType(tc.in))
		})
	}
}

func TestResourceComponent(t *testing.T) {
	var (
		tcc = []struct {
			in  string
			exp string
		}{
			{"ns::cmp:r/1/2/3", "ns::cmp"},
			{"ns::cmp:r", "ns::cmp"},
			{"ns::cmp/", "ns::cmp"},
			{"ns::cmp", "ns::cmp"},
			{"cmp", "cmp"},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.in, func(t *testing.T) {
			require.Equal(t, tc.exp, ResourceComponent(tc.in))
		})
	}
}

func TestResourceMatch(t *testing.T) {
	var (
		tcc = []struct {
			m string
			r string
			e bool
		}{
			{"::corteza:test/a/b/c", "::corteza:test/a/b/c", true},
			{"::corteza:test/a/b/*", "::corteza:test/a/b/c", true},
			{"::corteza:test/a/*/*", "::corteza:test/a/b/c", true},
			{"::corteza:test/*/*/*", "::corteza:test/a/b/c", true},
			{"::corteza:test/a/*/*", "::corteza:test/1/2/3", false},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.m, func(t *testing.T) {
			require.Equal(t, tc.e, matchResource(tc.m, tc.r))
		})
	}
}

//cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
//Benchmark_MatchResource100-16        	 7353837	       157.0 ns/op
//Benchmark_MatchResource1000-16       	 6868928	       166.0 ns/op
//Benchmark_MatchResource10000-16      	 7373701	       164.8 ns/op
//Benchmark_MatchResource100000-16     	 7556944	       156.5 ns/op
//Benchmark_MatchResource1000000-16    	 7445456	       157.8 ns/op
func benchmarkMatchResource(b *testing.B, c int) {
	b.StartTimer()

	for n := 0; n < b.N; n++ {
		for i := 0; i < c; i++ {
			matchResource("corteza::test/a/1/1/1", "corteza::test/a/1/1/1")
			matchResource("corteza::test/a/*/*/1", "corteza::test/a/1/1/1")
		}
	}

	b.StopTimer()
}

func Benchmark_MatchResource100(b *testing.B)     { benchmarkMatchResource(b, 100) }
func Benchmark_MatchResource1000(b *testing.B)    { benchmarkMatchResource(b, 1000) }
func Benchmark_MatchResource10000(b *testing.B)   { benchmarkMatchResource(b, 10000) }
func Benchmark_MatchResource100000(b *testing.B)  { benchmarkMatchResource(b, 100000) }
func Benchmark_MatchResource1000000(b *testing.B) { benchmarkMatchResource(b, 1000000) }

func TestLevel(t *testing.T) {
	var (
		tcc = []struct {
			r string
			l int
		}{
			{"corteza::test/a/b/c", 111},
			{"corteza::test/a/b/*", 11},
			{"corteza::test/a/*/*", 1},
			{"corteza::test/*/*/*", 0},
			{"corteza::test/a/*/123", 101},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.r, func(t *testing.T) {
			require.Equal(t, tc.l, level(tc.r))
		})
	}
}
