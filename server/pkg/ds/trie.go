package ds

import "sync"

type (
	trieNode[K comparable, V any] struct {
		key      K
		count    int
		children map[K]*trieNode[K, V]

		lief  int
		value V
	}

	Trie[K comparable, V any] struct {
		Size int
		root *trieNode[K, V]

		mu sync.RWMutex
	}
)

func NewTrie[K comparable, V any]() *Trie[K, V] {
	return &Trie[K, V]{
		root: &trieNode[K, V]{
			children: make(map[K]*trieNode[K, V]),
		},
	}
}

func TrieUpsert[K comparable, V any](t *Trie[K, V], merge func(a, b V) V, value V, key ...K) {
	aux, ok := TrieSearch[K, V](t, key...)
	if ok {
		aux = merge(aux, value)
	} else {
		aux = value
	}

	TrieInsert[K, V](t, aux, key...)
}

func TrieInsert[K comparable, V any](t *Trie[K, V], value V, key ...K) {
	if len(key) == 0 {
		return
	}

	inc := true

	// If we have an exact match, then we don't increment the thing
	if _, ok := TrieSearch[K, V](t, key...); ok {
		inc = false
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, k := range key {
		if node.children[k] == nil {
			node.children[k] = &trieNode[K, V]{
				key:      k,
				count:    0,
				children: make(map[K]*trieNode[K, V], 4),
			}
		}

		if inc {
			node.children[k].count++
		}

		node = node.children[k]
	}

	node.value = value
	node.lief++
	t.Size++
}

func TrieSearch[K comparable, V any](t *Trie[K, V], key ...K) (v V, ok bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	node := t.root
	for _, k := range key {
		if node == nil {
			ok = false
			return
		}

		node = node.children[k]
	}

	if node == nil {
		ok = false
		return
	}

	if node.lief == 0 {
		ok = false
		return
	}

	return node.value, true
}

func TrieRemove[K comparable, V any](t *Trie[K, V], key ...K) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node := t.root
	for _, k := range key {
		if node.children[k] == nil {
			return
		}

		node.children[k].count--
		if node.children[k].count == 0 {
			delete(node.children, k)
			return
		}

		node = node.children[k]
		if node == nil {
			return
		}
	}

	if node.lief == 0 {
		return
	}

	node.lief--
	t.Size--
}
