package rbac

import (
	"context"
	"sort"
	"sync"
	"time"
)

type (
	usageCounter struct {
		index map[uint64]uint

		lock sync.RWMutex

		sigThreshold uint

		incChan chan uint64
		sigChan chan counterEntry
	}

	counterEntry struct {
		key   uint64
		count uint
	}

	MinHeap []counterEntry
)

func (svc *usageCounter) worstPerformers(n int) (out []uint64) {
	svc.lock.RLock()
	defer svc.lock.RUnlock()

	// Code to get n elements with the smallest count

	hh := make(MinHeap, 0, len(svc.index))
	for k, v := range svc.index {
		hh = append(hh, counterEntry{key: k, count: v})
	}

	sort.Sort(hh)

	for _, x := range hh {
		out = append(out, x.key)

		if len(out) >= n {
			return
		}
	}

	return
}

func (svc *usageCounter) inc(key uint64) {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	count := svc.index[key] + 1
	svc.index[key] = count

	if count >= svc.sigThreshold {
		delete(svc.index, key)
		svc.sigChan <- counterEntry{key: key, count: count}
	}
}

func (svc *usageCounter) clean() {
	svc.lock.Lock()
	defer svc.lock.Unlock()

	for k, v := range svc.index {
		if v < uint(float64(svc.sigThreshold)*0.05) {
			delete(svc.index, k)
		}
	}
}

func (svc *usageCounter) watch(ctx context.Context) {
	cleanT := time.NewTicker(time.Minute * 10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case <-cleanT.C:
				svc.clean()

			case key := <-svc.incChan:
				svc.inc(key)
			}
		}
	}()
}

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
