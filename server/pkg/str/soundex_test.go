package str

import (
	"testing"
)

func Test_soundex(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			"Robert",
			"R163",
		},
		{
			"Rupert",
			"R163",
		},
		{
			"Rubin",
			"R150",
		},
		{
			"Ashcraft",
			"A261",
		},
		{
			"Ashcroft",
			"A261",
		},
		{
			"Tymczak",
			"T522",
		},
		{
			"Pfister",
			"P123",
		},
		{
			"AH KEY",
			"A000",
		},
		{
			"The quick brown fox",
			"T221",
		},
		{
			"h3110 w021d",
			"3000",
		},
		{
			"1337",
			"1000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSoundex(tt.name); got != tt.want {
				t.Errorf("soundex() = %v, want %v", got, tt.want)
			}
		})
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/str
// BenchmarkSoundex_100_100-12               475930              2835 ns/op            3568 B/op         80 allocs/op
// BenchmarkSoundex_1000_1000-12              30688             37494 ns/op          225458 B/op        657 allocs/op
// BenchmarkSoundex_10000_10000-12              568           2023842 ns/op        25114459 B/op       6892 allocs/op
// BenchmarkSoundex_100000_100000-12              7         160967804 ns/op        2544780057 B/op    69015 allocs/op
func benchmarkSoundex(b *testing.B, w1l int) {
	w1 := randStringRunes(w1l)

	for i := 0; i < b.N; i++ {
		ToSoundex(w1)
	}
}

func BenchmarkSoundex_100_100(b *testing.B) {
	benchmarkSoundex(b, 100)
}

func BenchmarkSoundex_1000_1000(b *testing.B) {
	benchmarkSoundex(b, 1000)
}

func BenchmarkSoundex_10000_10000(b *testing.B) {
	benchmarkSoundex(b, 10000)
}

func BenchmarkSoundex_100000_100000(b *testing.B) {
	benchmarkSoundex(b, 100000)
}
