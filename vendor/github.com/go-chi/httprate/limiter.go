package httprate

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

type LimitCounter interface {
	Increment(key string, currentWindow time.Time) error
	Get(key string, previousWindow, currentWindow time.Time) (int, int, error)
}

func NewRateLimiter(requestLimit int, windowLength time.Duration, counter LimitCounter, keyFuncs ...KeyFunc) *rateLimiter {
	var keyFn KeyFunc
	if len(keyFuncs) == 0 {
		keyFn = func(r *http.Request) (string, error) {
			return "*", nil
		}
	} else {
		keyFn = composedKeyFunc(keyFuncs...)
	}

	if counter == nil {
		counter = &localCounter{
			counters:     make(map[uint64]*count),
			windowLength: windowLength,
		}
	}

	return &rateLimiter{
		requestLimit: requestLimit,
		windowLength: windowLength,
		keyFn:        keyFn,
		limitCounter: counter,
	}
}

func LimitCounterKey(key string, window time.Time) uint64 {
	h := xxhash.New()
	h.WriteString(key)
	h.WriteString(fmt.Sprintf("%d", window.Unix()))
	return h.Sum64()
}

type rateLimiter struct {
	requestLimit int
	windowLength time.Duration
	keyFn        KeyFunc
	limitCounter LimitCounter
}

func (r *rateLimiter) Counter() LimitCounter {
	return r.limitCounter
}

func (r *rateLimiter) Status(key string) (bool, float64, error) {
	t := time.Now().UTC()
	currentWindow := t.Truncate(r.windowLength)
	previousWindow := currentWindow.Add(-r.windowLength)

	currCount, prevCount, err := r.limitCounter.Get(key, currentWindow, previousWindow)
	if err != nil {
		return false, 0, err
	}

	diff := t.Sub(currentWindow)
	rate := float64(prevCount)*(float64(r.windowLength)-float64(diff))/float64(r.windowLength) + float64(currCount)

	if rate > float64(r.requestLimit) {
		return false, rate, nil
	}
	return true, rate, nil
}

func (l *rateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := l.keyFn(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionRequired)
			return
		}

		currentWindow := time.Now().UTC().Truncate(l.windowLength)

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", l.requestLimit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", 0))
		w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", currentWindow.Add(l.windowLength).Unix()))

		_, rate, err := l.Status(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusPreconditionRequired)
			return
		}
		nrate := int(math.Round(rate))

		if l.requestLimit > nrate {
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", l.requestLimit-nrate))
		}

		if nrate >= l.requestLimit {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		err = l.limitCounter.Increment(key, currentWindow)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type localCounter struct {
	counters     map[uint64]*count
	windowLength time.Duration
	lastEvict    time.Time
	mu           sync.Mutex
}

var _ LimitCounter = &localCounter{}

type count struct {
	value     int
	updatedAt time.Time
}

func (c *localCounter) Increment(key string, currentWindow time.Time) error {
	c.evict()

	c.mu.Lock()
	defer c.mu.Unlock()

	hkey := LimitCounterKey(key, currentWindow)

	v, ok := c.counters[hkey]
	if !ok {
		v = &count{}
		c.counters[hkey] = v
	}
	v.value += 1
	v.updatedAt = time.Now()

	return nil
}

func (c *localCounter) Get(key string, currentWindow, previousWindow time.Time) (int, int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	curr, ok := c.counters[LimitCounterKey(key, currentWindow)]
	if !ok {
		curr = &count{value: 0, updatedAt: time.Now()}
	}
	prev, ok := c.counters[LimitCounterKey(key, previousWindow)]
	if !ok {
		prev = &count{value: 0, updatedAt: time.Now()}
	}

	return curr.value, prev.value, nil
}

func (c *localCounter) evict() {
	c.mu.Lock()
	defer c.mu.Unlock()

	d := c.windowLength * 3

	if time.Since(c.lastEvict) < d {
		return
	}

	for k, v := range c.counters {
		if time.Since(v.updatedAt) >= d {
			delete(c.counters, k)
		}
	}
}
