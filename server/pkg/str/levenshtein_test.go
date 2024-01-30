package str

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want int
	}{
		{"", "hello", 5},
		{"hello", "", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "ba", 2},
		{"ab", "aaa", 2},
		{"bbb", "a", 3},
		{"kitten", "sitting", 3},
		{"distance", "difference", 5},
		{"levenshtein", "frankenstein", 6},
		{"resume and cafe", "resumes and cafes", 2},
		{"a very long string that is meant to exceed", "another very long string that is meant to exceed", 6},
		// Testing acutes and umlauts
		{"resumé and café", "resumés and cafés", 2},
		{"resume and cafe", "resumé and café", 4},
		{"Hafþór Júlíus Björnsson", "Hafþor Julius Bjornsson", 8},
		// Only 2 characters are less in the 2nd string
		{"།་གམ་འས་པ་་མ།", "།་གམའས་པ་་མ", 6},
	}
	for _, tt := range tests {
		t.Run(tt.a, func(t *testing.T) {
			if got := ToLevenshteinDistance(tt.a, tt.b); got != tt.want {
				t.Errorf("LevenshteinDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/str
// BenchmarkLeven_100_100-12                  39949             25767 ns/op           93184 B/op        102 allocs/op
// BenchmarkLeven_1000_1000-12                  390           3081967 ns/op         8298552 B/op       1011 allocs/op
// BenchmarkLeven_10000_10000-12                  4         299103531 ns/op        829957216 B/op     10131 allocs/op
// PASS

func benchmarkLeven(b *testing.B, w1l, w2l int) {
	w1 := randStringRunes(w1l)
	w2 := randStringRunes(w2l)

	for i := 0; i < b.N; i++ {
		ToLevenshteinDistance(w1, w2)
	}
}

func BenchmarkLeven_100_100(b *testing.B) {
	benchmarkLeven(b, 100, 100)
}

func BenchmarkLeven_1000_1000(b *testing.B) {
	benchmarkLeven(b, 1000, 1000)
}

func BenchmarkLeven_10000_10000(b *testing.B) {
	benchmarkLeven(b, 10000, 10000)
}
