package types

import (
	"math/rand"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/ds"
	"github.com/stretchr/testify/require"
)

// pre
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/compose/types
// BenchmarkModuleFieldRbacResource-12    	 4533682	       253.2 ns/op	     204 B/op	       9 allocs/op

// post
// goos: darwin
// goarch: arm64
// pkg: github.com/cortezaproject/corteza/server/compose/types
// BenchmarkModuleFieldRbacResource-12    	26023652	        46.57 ns/op	       0 B/op	       0 allocs/op
func BenchmarkModuleFieldRbacResource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ModuleFieldRbacResource(1000+randomN(2), 2000+randomN(4), 3000+randomN(1000))
	}
}

func randomN(mx int) uint64 {
	return uint64(rand.Intn(mx) + 1)
}

func TestResourceIndex(t *testing.T) {
	resourceIndexMaxSize = 2

	const (
		namespaceID   uint64 = 1001
		moduleID      uint64 = 1002
		moduleFieldID uint64 = 1003
		pageID        uint64 = 1004
	)

	var (
		res string
		w   *indexWrapper
		ok  bool

		req = require.New(t)
	)

	// Access first res.

	res = ModuleRbacResource(namespaceID, moduleID)

	w, ok = ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, moduleID)
	req.Contains(res, "1001")
	req.Contains(res, "1002")
	req.Equal(1, resourceIndex.Size)
	req.True(ok)
	req.Equal(uint(1), w.counter)

	// Re-access the first res.

	res = ModuleRbacResource(namespaceID, moduleID)

	w, ok = ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, moduleID)
	req.Contains(res, "1001")
	req.Contains(res, "1002")
	req.Equal(1, resourceIndex.Size)
	req.True(ok)
	req.Equal(uint(2), w.counter)

	// Access second res.

	res = PageRbacResource(namespaceID, pageID)

	w, ok = ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, pageID)
	req.Contains(res, "1001")
	req.Contains(res, "1004")
	req.Equal(2, resourceIndex.Size)
	req.True(ok)
	req.Equal(uint(1), w.counter)

	// Access third res.; index cleanup (for now, it's just reset)

	res = ModuleFieldRbacResource(namespaceID, moduleID, moduleFieldID)

	w, ok = ds.TrieSearch[uint64, *indexWrapper](resourceIndex, namespaceID, moduleID, moduleFieldID)
	req.Contains(res, "1001")
	req.Contains(res, "1002")
	req.Contains(res, "1003")
	req.Equal(1, resourceIndex.Size)
	req.True(ok)
	req.Equal(uint(1), w.counter)
}
