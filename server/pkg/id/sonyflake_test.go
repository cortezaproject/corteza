package id

import (
	"testing"

	"github.com/crusttech/human/server/pkg/cli"
)

// goos: darwin
// goarch: arm64
// pkg: github.com/crusttech/human/server/pkg/id
// BenchmarkGenerator-12    	  162234	     39011 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGenerator(b *testing.B) {
	ctx := cli.Context()
	Init(ctx)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Next()
	}
}
