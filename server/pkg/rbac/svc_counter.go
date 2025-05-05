package rbac

import (
	"context"
	"sort"
	"sync"
	"time"
)

type (
	usageCounter[K comparable] struct {
		lock sync.RWMutex

		// index keeps track of all the things we're counting
		index map[K]counterItem[K]

		// sigEvictThreshold denotes when the usage counter should evict an item
		sigEvictThreshold float64
		// decayFactor denotes how fast the score decays
		// when 1 - it won't decay
		// when 0 - it's barely preserved
		decayFactor float64

		// incChan sends instructions to the counter re. key K increment
		incChan chan K

		// rmChan sends instructions to the counter re. key K removal
		rmChan chan uint64

		// checkKeyInclusion helps determine if an indexed thing should matches something
		checkKeyInclusion func(k K, role uint64) bool

		// decayInterval denotes in what interval the decay factor should apply
		decayInterval time.Duration
		// cleanupInterval denotes in what interval counter evicts stuff
		cleanupInterval time.Duration
	}

	// counterItem wraps some metadata around each index
	counterItem[K comparable] struct {
		key   K
		score float64

		// added denotes when the item was added to the counter
		added time.Time
		// lastScored denotes when the item was last scored (either via decay or access)
		lastScored time.Time
		// lastAccess denotes when the item was last accessed, needed
		lastAccess time.Time
	}

	MinHeap[K comparable] []counterItem[K]
)

// inc updates key K
func (svc *usageCounter[K]) inc(key K) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	_, ok := svc.index[key]
	if !ok {
		svc.procNew(key)
	} else {
		svc.procExisting(key)
	}
}

// evict evicts the items below the specified threshold
func (svc *usageCounter[K]) evict() (out []K) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	// Firstly score them up
	out = make([]K, 0, 4)
	for k, v := range svc.index {
		if v.score > float64(svc.sigEvictThreshold) {
			continue
		}

		out = append(out, k)
	}

	// Then delete them
	for _, r := range out {
		delete(svc.index, r)
	}

	return out
}

func (svc *usageCounter[K]) cleanRoleKeys(role uint64) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	for k := range svc.index {
		if !svc.checkKeyInclusion(k, role) {
			continue
		}

		delete(svc.index, k)
	}
}

// decay applies the specified decay factor to the cache items
func (svc *usageCounter[K]) decay() {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	n := time.Now()
	for k, v := range svc.index {
		if n.Before(v.lastAccess.Add(svc.decayInterval)) {
			continue
		}

		v.score *= svc.decayFactor
		svc.index[k] = v
	}
}

// bestPerformers returns the top n items based on their score
func (svc *usageCounter[K]) bestPerformers(n int) (out []K) {
	svc.lock.RLock()
	defer svc.lock.RUnlock()

	hh := make(MinHeap[K], 0, len(svc.index))
	for k, v := range svc.index {
		hh = append(hh, counterItem[K]{key: k, score: v.score})
	}

	sort.Sort(hh)

	for i := len(hh) - 1; i >= 0; i-- {
		if n >= 0 && len(out) >= n {
			return
		}

		out = append(out, hh[i].key)
	}
	return
}

// procNew notes a new key in the thing, defaults and stuff
func (svc *usageCounter[K]) procNew(key K) {
	n := time.Now()
	if svc.index == nil {
		svc.index = make(map[K]counterItem[K])
	}

	svc.index[key] = counterItem[K]{
		score:      1,
		added:      n,
		lastScored: n,
		lastAccess: n,
	}
}

// procExisting notes an access to an existing index element
func (svc *usageCounter[K]) procExisting(key K) {
	n := time.Now()

	aux := svc.index[key]
	aux.lastAccess = n
	aux.lastScored = n
	aux.score++

	svc.index[key] = aux
}

func (svc *usageCounter[K]) watch(ctx context.Context) {
	if svc.decayInterval == 0 {
		panic("svc.decayInterval can not be 0")
	}

	if svc.cleanupInterval == 0 {
		panic("svc.cleanupInterval can not be 0")
	}

	decayT := time.NewTicker(svc.decayInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-decayT.C:
				svc.decay()

			case key := <-svc.incChan:
				svc.inc(key)

			case role := <-svc.rmChan:
				svc.cleanRoleKeys(role)
			}
		}
	}()
}

func (h MinHeap[K]) Len() int           { return len(h) }
func (h MinHeap[K]) Less(i, j int) bool { return h[i].score < h[j].score }
func (h MinHeap[K]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
