package rbac

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/slice"
	"go.uber.org/zap"
)

type (
	StatsLogger struct {
		lock sync.RWMutex
		log  *zap.Logger

		// Channels for async comms
		cacheHitChan       chan statsWrap
		cacheMissChan      chan statsWrap
		timingDatabaseChan chan time.Duration
		timingIndexChan    chan time.Duration

		// Counters
		cacheHits         uint
		cacheMisses       uint
		cacheUpdates      uint
		avgDatabaseTiming time.Duration
		minDatabaseTiming time.Duration
		maxDatabaseTiming time.Duration
		avgIndexTiming    time.Duration
		minIndexTiming    time.Duration
		maxIndexTiming    time.Duration

		// Track a limited set of things
		// Using a circular buffer we can easily not consume too much data
		lastHits            *slice.Circular[string]
		lastMisses          *slice.Circular[string]
		lastDatabaseTimings *slice.Circular[time.Duration]
		lastIndexTimings    *slice.Circular[time.Duration]
	}

	// statsWrap wraps the state to log
	statsWrap struct {
		roles    []uint64
		resource string
		op       string
	}
)

const (
	maxLastLogs = 500
)

// Stats returns the tracked stats
func (l *StatsLogger) Stats() (
	cacheHit uint,
	cacheMiss uint,
	cacheUpdates uint,
	avgDbTiming, minDbTiming, maxDbTiming time.Duration,
	avgIndexTiming, minIndexTiming, maxIndexTiming time.Duration,
	lastHits []string,
	lastMisses []string,
	lastDbTimings []time.Duration,
	lastIndexTimings []time.Duration,
) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.cacheHits,
		l.cacheMisses,
		l.cacheUpdates,
		l.avgDatabaseTiming,
		l.minDatabaseTiming,
		l.maxDatabaseTiming,
		l.avgIndexTiming,
		l.minIndexTiming,
		l.maxIndexTiming,
		l.lastHits.Slice(),
		l.lastMisses.Slice(),
		l.lastDatabaseTimings.Slice(),
		l.lastIndexTimings.Slice()
}

// TimingDatabase logs the giving duration
func (l *StatsLogger) TimingDatabase(timing time.Duration) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.log.Info("record database timing", zap.Duration("timing", timing))

	{
		l.avgDatabaseTiming = (l.avgDatabaseTiming + timing) / 2
	}

	{
		if l.minDatabaseTiming == 0 {
			l.minDatabaseTiming = timing
		}
		if timing < l.minDatabaseTiming {
			l.minDatabaseTiming = timing
		}
	}

	{
		if l.maxDatabaseTiming == 0 {
			l.maxDatabaseTiming = timing
		}
		if timing > l.maxDatabaseTiming {
			l.maxDatabaseTiming = timing
		}
	}

	{
		if l.lastDatabaseTimings == nil {
			l.lastDatabaseTimings = slice.NewCircular[time.Duration](maxLastLogs)
		}

		l.lastDatabaseTimings.Add(timing)
	}
}

// TimingIndex logs the giving duration
func (l *StatsLogger) TimingIndex(timing time.Duration) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.log.Info("record index timing", zap.Duration("timing", timing))

	{
		l.avgIndexTiming = (l.avgIndexTiming + timing) / 2
	}

	{
		if l.minIndexTiming == 0 {
			l.minIndexTiming = timing
		}
		if timing < l.minIndexTiming {
			l.minIndexTiming = timing
		}
	}

	{
		if l.maxIndexTiming == 0 {
			l.maxIndexTiming = timing
		}
		if timing > l.maxIndexTiming {
			l.maxIndexTiming = timing
		}
	}

	{
		if l.lastIndexTimings == nil {
			l.lastIndexTimings = slice.NewCircular[time.Duration](maxLastLogs)
		}

		l.lastIndexTimings.Add(timing)
	}
}

func (l *StatsLogger) CacheHit(roles []uint64, resource string, op string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.log.Info("cache hit", zap.Any("roles", roles), zap.String("resource", resource), zap.String("op", op))

	l.cacheHits++
	if l.lastHits == nil {
		l.lastHits = slice.NewCircular[string](maxLastLogs)
	}
	l.lastHits.Add(l.strfEntry(roles, resource, op))
}

func (l *StatsLogger) CacheMiss(roles []uint64, resource string, op string) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.log.Info("cache miss", zap.Any("roles", roles), zap.String("resource", resource), zap.String("op", op))

	l.cacheMisses++
	if l.lastMisses == nil {
		l.lastMisses = slice.NewCircular[string](maxLastLogs)
	}
	l.lastMisses.Add(l.strfEntry(roles, resource, op))
}

func (l *StatsLogger) CacheUpdate(in *Rule) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.log.Info("cache update", zap.Any("rule", in))

	l.cacheUpdates++
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utils

func (l *StatsLogger) strfEntry(roles []uint64, resource string, op string) string {
	sort.Slice(roles, func(i, j int) bool { return roles[i] < roles[j] })

	return fmt.Sprintf("%v %s %s", roles, op, resource)
}

func (l *StatsLogger) watch(ctx context.Context) {
	t := time.NewTicker(time.Minute * 5)

	go func() {
		for {
			select {
			case <-t.C:
				l.log.Info("stats logger tick")

			case rs := <-l.cacheMissChan:
				l.CacheMiss(rs.roles, rs.resource, rs.op)

			case rs := <-l.cacheHitChan:
				l.CacheHit(rs.roles, rs.resource, rs.op)

			case tt := <-l.timingDatabaseChan:
				l.TimingDatabase(tt)

			case tt := <-l.timingIndexChan:
				l.TimingIndex(tt)

			case <-ctx.Done():
				l.log.Info("terminating watcher")
			}
		}
	}()
}
