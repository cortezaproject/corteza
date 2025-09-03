package id

import (
	"encoding/json"
	"testing"
)

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/id
// cpu: Apple M3 Pro
// BenchmarkIDCastingBase-12    	 2561474	       465.7 ns/op	     312 B/op	       8 allocs/op
// PASS
func BenchmarkIDCastingBase(b *testing.B) {
	b.ReportAllocs()

	type TestStruct struct {
		K1, K2 string
	}

	aux := TestStruct{
		K1: "K1 Value",
		K2: "K2 Value",
	}
	var dst TestStruct

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bb, err := json.Marshal(aux)
		if err != nil {
			b.Fatal(err)
		}
		if err := json.Unmarshal(bb, &dst); err != nil {
			b.Fatal(err)
		}
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/id
// cpu: Apple M3 Pro
// BenchmarkIDCasting-12    	 1284855	       866.6 ns/op	     528 B/op	      12 allocs/op
// PASS
func BenchmarkIDCasting(b *testing.B) {
	b.ReportAllocs()

	type TestStruct struct {
		SomethingID ID
		K1, K2      string
	}

	aux := TestStruct{
		SomethingID: MustNumID(3617289),
		K1:          "K1 Value",
		K2:          "K2 Value",
	}
	var dst TestStruct

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		bb, err := json.Marshal(aux)
		if err != nil {
			b.Fatal(err)
		}
		if err := json.Unmarshal(bb, &dst); err != nil {
			b.Fatal(err)
		}
	}
}
