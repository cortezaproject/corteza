package dal

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func benchmarkExecAggregate(b *testing.B, n int) {
	b.StopTimer()
	ctx := context.Background()
	group := []simpleAttribute{{ident: "g", expr: "g"}}
	outAttributes := []simpleAttribute{{ident: "sum", expr: "sum(s1_v)"}}
	sourceAttributes := []simpleAttribute{
		{ident: "k"},
		{ident: "g"},
		{ident: "s1_v"},
	}

	ggs := []string{"a", "b", "c", "d", "e", "f", "g"}

	// Inmem buffer for example
	buff := InMemoryBuffer()
	for i := 0; i < n; i++ {
		require.NoError(b, buff.Add(ctx, simpleRow{"k": i + 1, "g": ggs[rand.Intn(len(ggs))], "s1_v": rand.Intn(200)}))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		def := Aggregate{
			Ident:            "agg",
			Group:            saToMapping(group...),
			OutAttributes:    saToMapping(outAttributes...),
			SourceAttributes: saToMapping(sourceAttributes...),
		}

		b.StartTimer()
		_, err := def.iterator(ctx, buff)
		require.NoError(b, err)
		b.StopTimer()

		buff.Seek(ctx, 0)
	}
}

// goos: linux
// goarch: amd64
// pkg: github.com/cortezaproject/corteza-server/pkg/dal
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkExecAggregate_20000-12     	      63	  17902131 ns/op	 1423343 B/op	   80432 allocs/op
// BenchmarkExecAggregate_40000-12     	      33	  36179481 ns/op	 2783435 B/op	  160432 allocs/op
// BenchmarkExecAggregate_60000-12     	      21	  53681403 ns/op	 4143471 B/op	  240433 allocs/op
// BenchmarkExecAggregate_80000-12     	      15	  72234005 ns/op	 5503370 B/op	  320432 allocs/op
// BenchmarkExecAggregate_100000-12    	      12	  92145195 ns/op	 6863674 B/op	  400434 allocs/op
// BenchmarkExecAggregate_120000-12    	      10	 109842507 ns/op	 8223711 B/op	  480433 allocs/op
// BenchmarkExecAggregate_140000-12    	       8	 129250726 ns/op	 9583616 B/op	  560433 allocs/op
// BenchmarkExecAggregate_160000-12    	       7	 144896012 ns/op	10943776 B/op	  640434 allocs/op
// BenchmarkExecAggregate_180000-12    	       7	 162705594 ns/op	12303718 B/op	  720434 allocs/op
// BenchmarkExecAggregate_200000-12    	       6	 181970479 ns/op	13663848 B/op	  800434 allocs/op
// BenchmarkExecAggregate_220000-12    	       5	 204168375 ns/op	15024054 B/op	  880435 allocs/op
// BenchmarkExecAggregate_240000-12    	       5	 218573247 ns/op	16383470 B/op	  960432 allocs/op
// BenchmarkExecAggregate_260000-12    	       5	 238612939 ns/op	17744016 B/op	 1040435 allocs/op
// BenchmarkExecAggregate_280000-12    	       4	 261442973 ns/op	19103558 B/op	 1120433 allocs/op
// BenchmarkExecAggregate_300000-12    	       4	 342276668 ns/op	20464244 B/op	 1200436 allocs/op

func BenchmarkExecAggregate_20000(b *testing.B)  { benchmarkExecAggregate(b, 20000) }
func BenchmarkExecAggregate_40000(b *testing.B)  { benchmarkExecAggregate(b, 40000) }
func BenchmarkExecAggregate_60000(b *testing.B)  { benchmarkExecAggregate(b, 60000) }
func BenchmarkExecAggregate_80000(b *testing.B)  { benchmarkExecAggregate(b, 80000) }
func BenchmarkExecAggregate_100000(b *testing.B) { benchmarkExecAggregate(b, 100000) }
func BenchmarkExecAggregate_120000(b *testing.B) { benchmarkExecAggregate(b, 120000) }
func BenchmarkExecAggregate_140000(b *testing.B) { benchmarkExecAggregate(b, 140000) }
func BenchmarkExecAggregate_160000(b *testing.B) { benchmarkExecAggregate(b, 160000) }
func BenchmarkExecAggregate_180000(b *testing.B) { benchmarkExecAggregate(b, 180000) }
func BenchmarkExecAggregate_200000(b *testing.B) { benchmarkExecAggregate(b, 200000) }
func BenchmarkExecAggregate_220000(b *testing.B) { benchmarkExecAggregate(b, 220000) }
func BenchmarkExecAggregate_240000(b *testing.B) { benchmarkExecAggregate(b, 240000) }
func BenchmarkExecAggregate_260000(b *testing.B) { benchmarkExecAggregate(b, 260000) }
func BenchmarkExecAggregate_280000(b *testing.B) { benchmarkExecAggregate(b, 280000) }
func BenchmarkExecAggregate_300000(b *testing.B) { benchmarkExecAggregate(b, 300000) }
