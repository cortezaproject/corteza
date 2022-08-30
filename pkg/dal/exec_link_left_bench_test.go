package dal

import (
	"context"
	"math/rand"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/stretchr/testify/require"
)

func benchmarkExecLink_left(b *testing.B, n int) {
	ctx := context.Background()
	lattrs := []simpleAttribute{
		{ident: "l_k"},
		{ident: "l_v1"},
		{ident: "l_v2"},
	}
	fattrs := []simpleAttribute{
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
		def := Link{
			Ident:              "link",
			OutLeftAttributes:  saToMapping(lattrs...),
			OutRightAttributes: saToMapping(fattrs...),
			LeftAttributes:     saToMapping(la...),
			RightAttributes:    saToMapping(fa...),
			filter:             internalFilter{orderBy: filter.SortExprSet{{Column: "f_k"}}},
			On:                 LinkPredicate{Left: "l_k", Right: "f_ref"},
		}

		def.init(ctx)
		def.exec(ctx, l, f)

		l.Seek(ctx, 0)
		f.Seek(ctx, 0)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza-server/pkg/dal
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkExecLink_left_200-12     	    1874	    633602 ns/op
// BenchmarkExecLink_left_400-12     	     910	   1369591 ns/op
// BenchmarkExecLink_left_600-12     	     590	   2162512 ns/op
// BenchmarkExecLink_left_800-12     	     451	   2725109 ns/op
// BenchmarkExecLink_left_1000-12    	     358	   3568684 ns/op
// BenchmarkExecLink_left_1200-12    	     260	   4059573 ns/op
// BenchmarkExecLink_left_1400-12    	     250	   4742408 ns/op
// BenchmarkExecLink_left_1600-12    	     219	   5383900 ns/op
// BenchmarkExecLink_left_1800-12    	     181	   6880839 ns/op
// BenchmarkExecLink_left_2000-12    	     165	   6821432 ns/op
// BenchmarkExecLink_left_2200-12    	     134	   8386884 ns/op
// BenchmarkExecLink_left_2400-12    	     142	   8471661 ns/op
// BenchmarkExecLink_left_2600-12    	     135	   8702328 ns/op
// BenchmarkExecLink_left_2800-12    	     126	   9375744 ns/op
// BenchmarkExecLink_left_3000-12    	     114	  10156339 ns/op
// BenchmarkExecLink_left_3200-12    	      96	  10682149 ns/op
// BenchmarkExecLink_left_3400-12    	     104	  11361907 ns/op
// BenchmarkExecLink_left_3600-12    	      80	  12616645 ns/op
// BenchmarkExecLink_left_3800-12    	      73	  13761230 ns/op
// BenchmarkExecLink_left_4000-12    	      81	  13795284 ns/op
// BenchmarkExecLink_left_4200-12    	      69	  15898381 ns/op
// BenchmarkExecLink_left_4400-12    	      72	  15162550 ns/op
// BenchmarkExecLink_left_4600-12    	      74	  15650527 ns/op
// BenchmarkExecLink_left_4800-12    	      70	  19573447 ns/op
// BenchmarkExecLink_left_5000-12    	      69	  17214960 ns/op
func BenchmarkExecLink_left_200(b *testing.B)  { benchmarkExecLink_left(b, 200) }
func BenchmarkExecLink_left_400(b *testing.B)  { benchmarkExecLink_left(b, 400) }
func BenchmarkExecLink_left_600(b *testing.B)  { benchmarkExecLink_left(b, 600) }
func BenchmarkExecLink_left_800(b *testing.B)  { benchmarkExecLink_left(b, 800) }
func BenchmarkExecLink_left_1000(b *testing.B) { benchmarkExecLink_left(b, 1000) }
func BenchmarkExecLink_left_1200(b *testing.B) { benchmarkExecLink_left(b, 1200) }
func BenchmarkExecLink_left_1400(b *testing.B) { benchmarkExecLink_left(b, 1400) }
func BenchmarkExecLink_left_1600(b *testing.B) { benchmarkExecLink_left(b, 1600) }
func BenchmarkExecLink_left_1800(b *testing.B) { benchmarkExecLink_left(b, 1800) }
func BenchmarkExecLink_left_2000(b *testing.B) { benchmarkExecLink_left(b, 2000) }
func BenchmarkExecLink_left_2200(b *testing.B) { benchmarkExecLink_left(b, 2200) }
func BenchmarkExecLink_left_2400(b *testing.B) { benchmarkExecLink_left(b, 2400) }
func BenchmarkExecLink_left_2600(b *testing.B) { benchmarkExecLink_left(b, 2600) }
func BenchmarkExecLink_left_2800(b *testing.B) { benchmarkExecLink_left(b, 2800) }
func BenchmarkExecLink_left_3000(b *testing.B) { benchmarkExecLink_left(b, 3000) }
func BenchmarkExecLink_left_3200(b *testing.B) { benchmarkExecLink_left(b, 3200) }
func BenchmarkExecLink_left_3400(b *testing.B) { benchmarkExecLink_left(b, 3400) }
func BenchmarkExecLink_left_3600(b *testing.B) { benchmarkExecLink_left(b, 3600) }
func BenchmarkExecLink_left_3800(b *testing.B) { benchmarkExecLink_left(b, 3800) }
func BenchmarkExecLink_left_4000(b *testing.B) { benchmarkExecLink_left(b, 4000) }
func BenchmarkExecLink_left_4200(b *testing.B) { benchmarkExecLink_left(b, 4200) }
func BenchmarkExecLink_left_4400(b *testing.B) { benchmarkExecLink_left(b, 4400) }
func BenchmarkExecLink_left_4600(b *testing.B) { benchmarkExecLink_left(b, 4600) }
func BenchmarkExecLink_left_4800(b *testing.B) { benchmarkExecLink_left(b, 4800) }
func BenchmarkExecLink_left_5000(b *testing.B) { benchmarkExecLink_left(b, 5000) }
