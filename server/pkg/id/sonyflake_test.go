package id

import (
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/cli"
)

func BenchmarkGenerator(b *testing.B) {
	ctx := cli.Context()
	Init(ctx)

	for n := 0; n < b.N; n++ {
		Next()
	}
}
