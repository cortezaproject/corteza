package ds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrie(t *testing.T) {
	tr := Trie[string, int]()
	req := require.New(t)

	TrieInsert[string, int](tr, 1, "l1 n1", "l2 n1", "l3 n1")
	TrieInsert[string, int](tr, 2, "l1 n1", "l2 n1", "l3 n2")
	TrieInsert[string, int](tr, 3, "l1 n1", "l2 n2", "l3 n1")
	TrieInsert[string, int](tr, 4, "l1 n1", "l2 n2", "l3 n1", "l4 n1")
	TrieInsert[string, int](tr, 4, "l1 n1", "l2 n2", "l3 n1", "l4 n1")

	req.Equal(1, must(TrieSearch[string, int](tr, "l1 n1", "l2 n1", "l3 n1")))
	req.Equal(2, must(TrieSearch[string, int](tr, "l1 n1", "l2 n1", "l3 n2")))
	req.Equal(3, must(TrieSearch[string, int](tr, "l1 n1", "l2 n2", "l3 n1")))
	req.Equal(4, must(TrieSearch[string, int](tr, "l1 n1", "l2 n2", "l3 n1", "l4 n1")))

	TrieRemove[string, int](tr, "l1 n1", "l2 n1", "l3 n1")
	req.False(mustNot(TrieSearch[string, int](tr, "l1 n1", "l2 n1", "l3 n1")))

	TrieRemove[string, int](tr, "l1 n1", "l2 n1", "l3 n2")
	req.False(mustNot(TrieSearch[string, int](tr, "l1 n1", "l2 n1", "l3 n2")))

	TrieRemove[string, int](tr, "l1 n1", "l2 n2", "l3 n1")
	req.False(mustNot(TrieSearch[string, int](tr, "l1 n1", "l2 n2", "l3 n1")))

	TrieRemove[string, int](tr, "l1 n1", "l2 n2", "l3 n1", "l4 n1")
	req.False(mustNot(TrieSearch[string, int](tr, "l1 n1", "l2 n2", "l3 n1", "l4 n1")))

	req.Len(tr.root.children, 0)
}

// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/pkg/ds
// BenchmarkTrie-12    	 4713706	       248.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTrie(b *testing.B) {
	tr := Trie[string, int]()

	bit1 := []string{"l1 n1", "l2 n1", "l3 n1"}
	bit2 := []string{"l1 n1", "l2 n1", "l3 n2"}
	bit3 := []string{"l1 n1", "l2 n2", "l3 n1"}
	bit4 := []string{"l1 n1", "l2 n2", "l3 n1", "l4 n1"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TrieInsert[string, int](tr, 1, bit1...)
		TrieInsert[string, int](tr, 2, bit2...)
		TrieInsert[string, int](tr, 3, bit3...)
		TrieInsert[string, int](tr, 4, bit4...)
		TrieInsert[string, int](tr, 4, bit4...)

		TrieSearch[string, int](tr, bit1...)
		TrieSearch[string, int](tr, bit2...)
		TrieSearch[string, int](tr, bit3...)
		TrieSearch[string, int](tr, bit4...)

		TrieSearch[string, int](tr, bit1...)
		TrieSearch[string, int](tr, bit2...)
		TrieSearch[string, int](tr, bit3...)
		TrieSearch[string, int](tr, bit4...)
	}
}

func must(v any, ok bool) any {
	if !ok {
		panic("must be found")
	}

	return v
}

func mustNot(_ any, ok bool) bool {
	if ok {
		panic("must not be found")
	}

	return false
}
