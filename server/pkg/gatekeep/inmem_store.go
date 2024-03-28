package gatekeep

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/cortezaproject/corteza/server/pkg/ds"
)

type (
	inmemStore struct {
		mux sync.RWMutex
		// We'll use the trie here so we can compress the data a bit
		// @todo consider switching to a plain'ol map as it might make more sense
		//       since the values would constantly be cleaned up
		index *ds.Trie[string, []byte]
	}
)

func InmemStore() *inmemStore {
	return &inmemStore{
		index: ds.NewTrie[string, []byte](),
	}
}

func (ix *inmemStore) GetValue(ctx context.Context, key string) (out []byte, err error) {
	ix.mux.RLock()
	defer ix.mux.RUnlock()

	out, ok := ds.TrieSearch[string, []byte](ix.index, strings.Split(key, "/")...)
	if !ok {
		err = fmt.Errorf("not found")
		return
	}

	return
}

func (ix *inmemStore) SetValue(ctx context.Context, key string, v []byte) error {
	ix.mux.Lock()
	defer ix.mux.Unlock()

	ds.TrieInsert[string, []byte](ix.index, v, strings.Split(key, "/")...)
	return nil
}

func (ix *inmemStore) DeleteValue(ctx context.Context, key string) error {
	ix.mux.Lock()
	defer ix.mux.Unlock()

	ds.TrieRemove[string, []byte](ix.index, strings.Split(key, "/")...)

	return nil
}
