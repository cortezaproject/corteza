package dal

import (
	"context"
	"math/rand"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func benchmarkExecJoin_local(b *testing.B, n int) {
	ctx := context.Background()
	attrs := []simpleAttribute{
		{ident: "l_k"},
		{ident: "l_v1"},
		{ident: "l_v2"},
		{ident: "f_k"},
		{ident: "f_ref"},
		{ident: "f_v1"},
		{ident: "f_v2"},
	}

	la := []simpleAttribute{
		{ident: "l_k", t: TypeID{}},
		{ident: "l_v1"},
		{ident: "l_v2"},
	}
	fa := []simpleAttribute{
		{ident: "f_k", t: TypeID{}},
		{ident: "f_ref", t: TypeID{}},
		{ident: "f_v1"},
		{ident: "f_v2"},
	}

	// Inmem buffer for example
	l := InMemoryBuffer()
	f := InMemoryBuffer()
	for i := 0; i < n; i++ {
		require.NoError(b, l.Add(ctx, simpleRow{"l_k": i + 1, "l_v1": "a", "l_v2": rand.Intn(200)}))
		require.NoError(b, f.Add(ctx, simpleRow{"f_k": i + 1, "f_ref": i + 1, "f_v1": "a", "f_v2": rand.Intn(200)}))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		def := Join{
			Ident:           "join",
			OutAttributes:   saToMapping(attrs...),
			LeftAttributes:  saToMapping(la...),
			RightAttributes: saToMapping(fa...),
			filter:          internalFilter{orderBy: filter.SortExprSet{{Column: "f_k"}}},
			On:              JoinPredicate{Left: "l_k", Right: "f_ref"},
		}

		def.iterator(ctx, l, f)

		l.Seek(ctx, 0)
		f.Seek(ctx, 0)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza/server/pkg/dal
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkExecJoin_local_200-12     	    1620	    718632 ns/op
// BenchmarkExecJoin_local_400-12     	     801	   1474321 ns/op
// BenchmarkExecJoin_local_600-12     	     504	   2344433 ns/op
// BenchmarkExecJoin_local_800-12     	     388	   3098124 ns/op
// BenchmarkExecJoin_local_1000-12    	     304	   3876453 ns/op
// BenchmarkExecJoin_local_1200-12    	     258	   4653349 ns/op
// BenchmarkExecJoin_local_1400-12    	     218	   5406631 ns/op
// BenchmarkExecJoin_local_1600-12    	     192	   6215687 ns/op
// BenchmarkExecJoin_local_1800-12    	     168	   7185540 ns/op
// BenchmarkExecJoin_local_2000-12    	     148	   8021597 ns/op
// BenchmarkExecJoin_local_2200-12    	     122	   9350466 ns/op
// BenchmarkExecJoin_local_2400-12    	     100	  10035371 ns/op
// BenchmarkExecJoin_local_2600-12    	     117	  10247223 ns/op
// BenchmarkExecJoin_local_2800-12    	     100	  11282374 ns/op
// BenchmarkExecJoin_local_3000-12    	      99	  12591232 ns/op
// BenchmarkExecJoin_local_3200-12    	      86	  13047472 ns/op
// BenchmarkExecJoin_local_3400-12    	      80	  13859590 ns/op
// BenchmarkExecJoin_local_3600-12    	      76	  14986832 ns/op
// BenchmarkExecJoin_local_3800-12    	      70	  15826467 ns/op
// BenchmarkExecJoin_local_4000-12    	      61	  16537483 ns/op
// BenchmarkExecJoin_local_4200-12    	      68	  17226297 ns/op
// BenchmarkExecJoin_local_4400-12    	      67	  17996956 ns/op
// BenchmarkExecJoin_local_4600-12    	      62	  18622211 ns/op
// BenchmarkExecJoin_local_4800-12    	      61	  19948914 ns/op
// BenchmarkExecJoin_local_5000-12    	      57	  20614964 ns/op
func BenchmarkExecJoin_local_200(b *testing.B)  { benchmarkExecJoin_local(b, 200) }
func BenchmarkExecJoin_local_400(b *testing.B)  { benchmarkExecJoin_local(b, 400) }
func BenchmarkExecJoin_local_600(b *testing.B)  { benchmarkExecJoin_local(b, 600) }
func BenchmarkExecJoin_local_800(b *testing.B)  { benchmarkExecJoin_local(b, 800) }
func BenchmarkExecJoin_local_1000(b *testing.B) { benchmarkExecJoin_local(b, 1000) }
func BenchmarkExecJoin_local_1200(b *testing.B) { benchmarkExecJoin_local(b, 1200) }
func BenchmarkExecJoin_local_1400(b *testing.B) { benchmarkExecJoin_local(b, 1400) }
func BenchmarkExecJoin_local_1600(b *testing.B) { benchmarkExecJoin_local(b, 1600) }
func BenchmarkExecJoin_local_1800(b *testing.B) { benchmarkExecJoin_local(b, 1800) }
func BenchmarkExecJoin_local_2000(b *testing.B) { benchmarkExecJoin_local(b, 2000) }
func BenchmarkExecJoin_local_2200(b *testing.B) { benchmarkExecJoin_local(b, 2200) }
func BenchmarkExecJoin_local_2400(b *testing.B) { benchmarkExecJoin_local(b, 2400) }
func BenchmarkExecJoin_local_2600(b *testing.B) { benchmarkExecJoin_local(b, 2600) }
func BenchmarkExecJoin_local_2800(b *testing.B) { benchmarkExecJoin_local(b, 2800) }
func BenchmarkExecJoin_local_3000(b *testing.B) { benchmarkExecJoin_local(b, 3000) }
func BenchmarkExecJoin_local_3200(b *testing.B) { benchmarkExecJoin_local(b, 3200) }
func BenchmarkExecJoin_local_3400(b *testing.B) { benchmarkExecJoin_local(b, 3400) }
func BenchmarkExecJoin_local_3600(b *testing.B) { benchmarkExecJoin_local(b, 3600) }
func BenchmarkExecJoin_local_3800(b *testing.B) { benchmarkExecJoin_local(b, 3800) }
func BenchmarkExecJoin_local_4000(b *testing.B) { benchmarkExecJoin_local(b, 4000) }
func BenchmarkExecJoin_local_4200(b *testing.B) { benchmarkExecJoin_local(b, 4200) }
func BenchmarkExecJoin_local_4400(b *testing.B) { benchmarkExecJoin_local(b, 4400) }
func BenchmarkExecJoin_local_4600(b *testing.B) { benchmarkExecJoin_local(b, 4600) }
func BenchmarkExecJoin_local_4800(b *testing.B) { benchmarkExecJoin_local(b, 4800) }
func BenchmarkExecJoin_local_5000(b *testing.B) { benchmarkExecJoin_local(b, 5000) }
